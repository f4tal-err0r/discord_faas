package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	client *minio.Client
}

func NewMinio() (*Minio, error) {
	endpoint := "discord-faas:9000"
	rfaasuser := os.Getenv("MINIO_ROOT_USER")
	rfaaspass := os.Getenv("MINIO_ROOT_PASSWORD")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(rfaasuser, rfaaspass, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating minio client: %v", err)
	}

	initBuckets := []string{"faas-artifacts", "faas-data"}

	createBucket := func(ctx context.Context, client *minio.Client, bucketName string) error {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			exists, err := client.BucketExists(ctx, bucketName)
			if err == nil && exists {
				return nil
			} else {
				return err
			}
		}

		return nil
	}

	for _, bucket := range initBuckets {
		err = createBucket(context.Background(), minioClient, bucket)
		if err != nil {
			return nil, fmt.Errorf("error creating bucket: %v", err)
		}
	}

	return &Minio{
		client: minioClient,
	}, nil
}

func (m *Minio) AddSrcArtifact(ctx context.Context, name string, data io.Reader, size int64) error {

	_, err := m.client.PutObject(ctx, "faas-data", "src/"+name, data, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *Minio) GetPresignedUrl(ctx context.Context, bucket string, cmdid string) (*url.URL, error) {

	url, err := m.client.PresignedPutObject(ctx, bucket, cmdid+".func", (30 * time.Minute))
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (m *Minio) GetSrcPath(ctx context.Context, name string) (string, error) {
	return url.JoinPath("s3://faas-data", "src", name)
}

func (m *Minio) DeleteSrcArtifact(ctx context.Context, name string) error {
	err := m.client.RemoveObject(ctx, "faas-data", "src/"+name, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
