package aws

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type AWS interface {
	DownloadObject(path string) (*AWSStorageObject, error)
	DownloadObjects() ([]*AWSStorageObject, error)
	UploadObject(filename string, data []byte) error
}

type AWSManager struct {
	bucketName string
	client     *s3.Client
}

type AWSStorageObject struct {
	Name        string
	ContentType string
	Body        []byte
}

func NewAWSManager(bucketName string, awsCfg aws.Config) *AWSManager {
	client := s3.NewFromConfig(awsCfg)

	return &AWSManager{
		bucketName: bucketName,
		client:     client,
	}
}

func (m *AWSManager) DownloadObject(path string) (*AWSStorageObject, error) {
	object, err := m.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return &AWSStorageObject{}, err
	}

	defer object.Body.Close()

	body, err := io.ReadAll(object.Body)

	return &AWSStorageObject{
		Name:        path,
		ContentType: *object.ContentType,
		Body:        body,
	}, err
}

func (m *AWSManager) DownloadObjects(paths []string) ([]*AWSStorageObject, error) {
	var outputObjects []*AWSStorageObject

	for _, path := range paths {
		outputObject, err := m.DownloadObject(path)
		if err != nil {
			return []*AWSStorageObject{}, err
		}
		outputObjects = append(outputObjects, outputObject)
	}

	return outputObjects, nil
}

func (m *AWSManager) UploadObject(path string, data []byte) error {
	body := bytes.NewReader(data)

	_, err := m.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(path),
		Body:   body,
	})

	return err
}
