package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	minio "github.com/minio/minio-go"
)

// S3 - A S3 implements FileSystem using the minio client
// allowing access to your S3 buckets and objects.
//
// Note that S3 will allow all access to files in your private
// buckets, If you have any sensitive information please make
// sure to not sure this project.
type S3 struct {
	*minio.Client
	bucket string
}

// Open - implements http.Filesystem implementation.
func (s3 *S3) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, "/") {
		return &httpMinioObject{
			client: s3.Client,
			object: nil,
			isDir:  true,
			bucket: bucket,
			prefix: name,
		}, nil
	}

	name = strings.TrimPrefix(name, "/")
	obj, err := s3.Client.GetObject(s3.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, os.ErrNotExist
	}

	if _, err = obj.Stat(); err != nil {
		return nil, os.ErrNotExist
	}

	return &httpMinioObject{
		client: s3.Client,
		object: obj,
		isDir:  false,
		bucket: bucket,
		prefix: name,
	}, nil
}

var (
	endpoint  string
	accessKey string
	secretKey string
	address   string
	bucket    string
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "https://s3.amazonaws.com", "S3 server endpoint.")
	flag.StringVar(&accessKey, "accessKey", "", "Access key of S3 storage.")
	flag.StringVar(&secretKey, "secretKey", "", "Secret key of S3 storage.")
	flag.StringVar(&bucket, "bucket", "", "Bucket name which hosts static files.")
	flag.StringVar(&address, "address", ":8080", "Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname.")
}

func main() {
	flag.Parse()

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := minio.NewV4(u.Host, accessKey, secretKey, u.Scheme == "https")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Started listening on", address)
	log.Fatalln(http.ListenAndServe(address, http.FileServer(&S3{client, bucket})))
}
