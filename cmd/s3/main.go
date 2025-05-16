package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

	ctx := context.Background()
	bucketName := os.Getenv("S3_BUCKET")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("LOCALSTACK_ENDPOINT"))
	})

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), *object.Size)
	}
}
