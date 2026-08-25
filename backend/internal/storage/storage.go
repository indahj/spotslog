package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client     *minio.Client
	bucket     string
	publicBase string
}

func New(endpoint, accessKey, secretKey, bucket string, useSSL bool, publicBase string) (*Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	s := &Storage{client: client, bucket: bucket, publicBase: publicBase}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	// Photo URLs are persisted in the database and rendered straight into
	// <img src="...">, so objects must be readable without a signature —
	// a fresh bucket is private and every photo would 404/403 in the browser.
	policy := fmt.Sprintf(`{
	  "Version": "2012-10-17",
	  "Statement": [{
	    "Effect": "Allow",
	    "Principal": {"AWS": ["*"]},
	    "Action": ["s3:GetObject"],
	    "Resource": ["arn:aws:s3:::%s/*"]
	  }]
	}`, bucket)
	if err := client.SetBucketPolicy(ctx, bucket, policy); err != nil {
		return nil, fmt.Errorf("set public-read policy on %q: %w", bucket, err)
	}

	return s, nil
}

// UploadPhoto uploads a multipart file under a namespaced key (e.g. "places/12/xyz.jpg")
// and returns the public URL to store alongside the record.
func (s *Storage) UploadPhoto(ctx context.Context, prefix string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s/%d-%s", prefix, time.Now().UnixNano(), fileHeader.Filename)

	_, err = s.client.PutObject(ctx, s.bucket, key, bytes.NewReader(buf), int64(len(buf)), minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.publicBase, key), nil
}
