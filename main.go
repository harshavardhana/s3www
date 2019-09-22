package main

import (
	"flag"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/minio/mc/pkg/console"
	minio "github.com/minio/minio-go/v6"
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
	obj, err := getObject(s3, name)
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

func getObject(s3 *S3, name string) (*minio.Object, error) {
	names := [3]string{name, name + "/index.html", name + "/index.htm"}
	for _, n := range names {
		obj, err := s3.Client.GetObject(s3.bucket, n, minio.GetObjectOptions{})

		if err == nil {
			if _, err = obj.Stat(); err == nil {
				return obj, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

var (
	endpoint  string
	accessKey string
	secretKey string
	address   string
	bucket    string
	tlsCert   string
	tlsKey    string
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "https://play.min.io", "S3 server endpoint.")
	flag.StringVar(&accessKey, "accessKey", "Q3AM3UQ867SPQQA43P2F", "Access key of S3 storage.")
	flag.StringVar(&secretKey, "secretKey", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "Secret key of S3 storage.")
	flag.StringVar(&bucket, "bucket", "testbucket", "Bucket name which hosts static files.")
	flag.StringVar(&address, "address", "127.0.0.1:8080", "Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname.")
	flag.StringVar(&tlsCert, "ssl-cert", "", "TLS certificate for this server.")
	flag.StringVar(&tlsKey, "ssl-key", "", "TLS private key for this server.")
}

func main() {
	flag.Parse()

	if strings.TrimSpace(bucket) == "" {
		console.Fatalln(`Bucket name cannot be empty, please provide 's3www -bucket "mybucket"'`)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		console.Fatalln(err)
	}

	client, err := minio.NewV4(u.Host, accessKey, secretKey, u.Scheme == "https")
	if err != nil {
		console.Fatalln(err)
	}

	if tlsCert != "" && tlsKey != "" {
		console.Infof("Started listening on https://%s\n", address)
		console.Fatalln(http.ListenAndServeTLS(address, tlsCert, tlsKey, http.FileServer(&S3{client, bucket})))
	} else {
		console.Infof("Started listening on http://%s\n", address)
		console.Fatalln(http.ListenAndServe(address, http.FileServer(&S3{client, bucket})))
	}
}
