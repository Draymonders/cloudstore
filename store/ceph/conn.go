package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"

	cfg "cloudstore/config"
)

var cephConn *s3.S3

// GetCephConn : get the ceph connection
func GetCephConn() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}
	// init ceph info
	auth := aws.Auth{
		AccessKey: cfg.CephAccessKey,
		SecretKey: cfg.CephSecretKey,
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          cfg.CephGWEndpoint,
		S3Endpoint:           cfg.CephGWEndpoint,
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}
	// create s3 connection
	return s3.New(auth, curRegion)
}

// GetCephBucket : get the name of bucket
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConn()
	return conn.Bucket(bucket)
}

// PutObject : put object to ceph (path, bucketname, datas)
func PutObject(bucketname string, path string, data []byte) error {
	return GetCephBucket(bucketname).Put(path, data, "octet-stream", s3.PublicRead)
}

// GetObject : get object from ceph
func GetObject(bucketname string, path string) ([]byte, error) {
	bucket := GetCephBucket(bucketname)
	fileData, err := bucket.Get(path)
	return fileData, err
}
