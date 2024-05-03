package util

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Bucket struct {
	S3Client       *s3.Client
	BucketName     string
	BucketLocation string
}

func NewS3Bucket(config aws.Config, bucketName string) *S3Bucket {
	return &S3Bucket{
		S3Client:   s3.NewFromConfig(config),
		BucketName: bucketName,
	}
}

func (s S3Bucket) UploadFile(objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(objectKey),
			Body:   file,
			ACL:    types.ObjectCannedACLPublicRead,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, s.BucketName, objectKey, err)
		}
	}
	return err
}

func (s S3Bucket) GetObjectUrl(objectKey string) (string, error) {
	location, err := s.S3Client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
		Bucket: aws.String(s.BucketName),
	})
	if err != nil {
		log.Printf("Couldn't get bucket location for %v. Here's why: %v\n", s.BucketName, err)
		return "", err
	}

	log.Printf("Got bucket location for %v: %v\n", s.BucketName, location.LocationConstraint)

	region := location.LocationConstraint
	if region == "" {
		region = types.BucketLocationConstraint("us-east-1")
	}

	url := "https://" + s.BucketName + ".s3." + string(region) + ".amazonaws.com/" + objectKey

	return url, nil
}
