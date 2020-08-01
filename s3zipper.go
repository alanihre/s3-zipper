package s3zipper

import (
	"archive/zip"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"sync"
)

func S3Zip(sess *session.Session, inputBucket string, files []string, archiveBucket string, archiveFileName string) error {
	// Create a buffer to write our archive to.
	pr, pw := io.Pipe()

	// Create a new zip archive.
	w := zip.NewWriter(pw)

	//Channel for passing errors from WaitGroup
	fatalErrors := make(chan error)
	//Channel for detecting when WaitGroup is finished
	wgDone := make(chan bool)

	// Wait for both download and upload to finish
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer pw.Close()

		downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
			// To enforce sequential download
			d.Concurrency = 1
		})

		for _, file := range files {
			// Create a file in the zip archive
			f, err := w.Create(file)
			if err != nil {
				fatalErrors <- err
				return
			}

			serialWriter := SerialWriterAt{
				w: f,
			}

			// Download file and write to zip archive via serialWriter
			_, err = downloader.Download(serialWriter, &s3.GetObjectInput{
				Bucket: aws.String(inputBucket),
				Key:    aws.String(file),
			})
			if err != nil {
				fatalErrors <- err
				return
			}
		}

		err := w.Close()
		if err != nil {
			fatalErrors <- err
			return
		}
	}()

	go func() {
		defer wg.Done()

		// Upload zip archive
		uploader := s3manager.NewUploader(sess)
		_, err := uploader.Upload(&s3manager.UploadInput{
			Body:        pr,
			Bucket:      aws.String(archiveBucket),
			ContentType: aws.String("application/zip"),
			Key:         aws.String(archiveFileName),
		})
		if err != nil {
			fatalErrors <- err
			return
		}
	}()

	// Final goroutine to wait until WaitGroup is done
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received through the channel
	select {
	case <-wgDone:
		// No errors
		return nil
	case err := <-fatalErrors:
		close(fatalErrors)
		return err
	}
}

type SerialWriterAt struct {
	w io.Writer
}

func (fw SerialWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	// Ignore 'offset' because we forced sequential downloads
	return fw.w.Write(p)
}
