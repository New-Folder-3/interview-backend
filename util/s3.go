package util

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
	"interview/internal/client"
	"interview/internal/conf"
	"io"
	"time"
)

var (
	s3Client *s3.Client
	s3Ready  bool
)

func S3Init() {
	if conf.Conf.S3.Use {
		if conf.Conf.S3.Path == "" ||
			conf.Conf.S3.AccessKey == "" ||
			conf.Conf.S3.SecretKey == "" ||
			conf.Conf.S3.Endpoint == "" ||
			conf.Conf.S3.Bucket == "" ||
			conf.Conf.S3.Region == "" {
			s3Ready = false
			return
		}
		awsCfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(conf.Conf.S3.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				conf.Conf.S3.AccessKey,
				conf.Conf.S3.SecretKey,
				"")),
			config.WithHTTPClient(client.GlobalHTTPClient),
		)
		if err != nil {
			ErrorPrinter(err)
			s3Ready = false
		}
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(conf.Conf.S3.Endpoint)
		})
		s3Ready = true
	}
}

func S3Upload(file io.ReadSeeker, fileName string) error {
	if !s3Ready {
		return errors.WithStack(errors.New("S3 Config or Init Failed"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &conf.Conf.S3.Bucket,
		Key:    &fileName,
		Body:   file,
		ACL:    types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func S3Download(fileName string) (io.ReadCloser, error) {
	if !s3Ready {
		return nil, errors.WithStack(errors.New("S3 Config or Init Failed"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	out, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &conf.Conf.S3.Bucket,
		Key:    &fileName,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return out.Body, nil
}
