package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	client *minio.Client
	opts   config.Storage
}

func NewS3Client(opts config.Storage) (*S3Client, error) {
	bucketName := opts.S3.Hostname
	s3Client, err := minio.New(bucketName, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.S3.Username, opts.S3.Password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating minio client: %v", err)
	}

	initBuckets := []string{bucketName}

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
		err = createBucket(context.Background(), s3Client, bucket)
		if err != nil {
			return nil, fmt.Errorf("error creating bucket: %v", err)
		}
	}

	return &S3Client{
		client: s3Client,
		opts:   opts,
	}, nil
}

func (m *S3Client) AddArtifact(ctx context.Context, name string, data io.Reader, size int64) error {
	_, err := m.client.PutObject(ctx, m.opts.S3.Hostname, name, data, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *S3Client) ListArtifacts(ctx context.Context, path string) ([]string, error) {
	objectCh := m.client.ListObjects(ctx, m.opts.S3.Hostname, minio.ListObjectsOptions{
		Recursive: true,
	})
	var objects []string
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object.Key)
	}
	return objects, nil
}

func (m *S3Client) GetArtifact(ctx context.Context, name string) (io.ReadCloser, error) {
	return m.client.GetObject(ctx, m.opts.S3.Hostname, name, minio.GetObjectOptions{})
}

func (m *S3Client) DeleteArtifact(ctx context.Context, name string) error {
	err := m.client.RemoveObject(ctx, m.opts.S3.Hostname, name, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
