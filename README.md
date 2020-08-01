# s3zipper
Go package for creating a zip archive of files on AWS S3 and upload the file to S3.

No disk space is used for this operation as the files are archived and piped directly to the zip archive file on S3 through a buffer. The full file is not stored in memory but is instead buffered to allow the application to run in limied memory environments.

## Example
```go
package main

import (
	"github.com/alanihre/s3zipper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	region := "eu-north-1"

	// Create a new AWS session. Credentials stored in environment variables
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}

	inputBucket := "my-s3-bucket"
	files := []string{"file1.txt", "file2.txt"}
	archiveBucket := "my-s3-bucket"
	archiveFileName := "archive.zip"

	err = s3zipper.S3Zip(sess, inputBucket, files, archiveBucket, archiveFileName)
	if err != nil {
		panic(err)
	}
}
```
