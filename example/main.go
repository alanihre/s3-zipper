package main

import (
	"github.com/alanihre/s3zipper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	// Create a new AWS session. Credentials stored in environment variables
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
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
