package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWS interface {
	DownloadObject(path string) ([]byte, error)
}

type AWSManager struct {
	bucketName string
	client     *s3.Client
}

func NewAWSManager(bucketName string, awsCfg aws.Config) *AWSManager {
	client := s3.NewFromConfig(awsCfg)

	return &AWSManager{
		bucketName: bucketName,
		client:     client,
	}
}

func (m *AWSManager) DownloadObject(path string) ([]byte, error) {
	// Getting object info
	objectInfo, err := m.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	// Creating an object to be returned
	body := make([]byte, objectInfo.ContentLength)
	w := manager.NewWriteAtBuffer(body)

	// Download file into the buffer
	downloader := manager.NewDownloader(m.client)
	_, err = downloader.Download(context.TODO(), w, &s3.GetObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(path),
	})

	return body, err
}
