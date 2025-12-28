package oss

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/Boyuan-IT-Club/go-kit/logs"
	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type OSSClient struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

func NewOSSClient(cfg *config.Config) *OSSClient {
	minioClient, err := minio.New(cfg.OSS.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.OSS.AccessKeyID, cfg.OSS.SecretAccessKey, ""),
		Secure: cfg.OSS.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.OSS.BucketName)
	if err != nil {
		logs.Errorf("Failed to check bucket existence: %v", err)
		panic(err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.OSS.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			logs.Errorf("Failed to create bucket: %v", err)
			panic(err)
		}
		
		// Set bucket policy to public read for download
		policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + cfg.OSS.BucketName + `/*"],"Sid": ""}]}`
		err = minioClient.SetBucketPolicy(ctx, cfg.OSS.BucketName, policy)
		if err != nil {
			logs.Errorf("Failed to set bucket policy: %v", err)
		}
	}

	return &OSSClient{
		client:     minioClient,
		bucketName: cfg.OSS.BucketName,
		endpoint:   cfg.OSS.Endpoint,
		useSSL:     cfg.OSS.UseSSL,
	}
}

func (c *OSSClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	_, err := c.client.PutObject(ctx, c.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Return public URL
	protocol := "http"
	if c.useSSL {
		protocol = "https"
	}
	return protocol + "://" + c.endpoint + "/" + c.bucketName + "/" + objectName, nil
}

func (c *OSSClient) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := c.client.PresignedGetObject(ctx, c.bucketName, objectName, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
