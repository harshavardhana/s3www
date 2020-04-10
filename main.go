package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/caddyserver/certmagic"
	minio "github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/credentials"
	"github.com/minio/minio-go/v6/pkg/s3utils"
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
	if strings.HasSuffix(name, pathSeparator) {
		return &httpMinioObject{
			client: s3.Client,
			object: nil,
			isDir:  true,
			bucket: bucket,
			prefix: strings.TrimSuffix(name, pathSeparator),
		}, nil
	}

	name = strings.TrimPrefix(name, pathSeparator)
	obj, err := getObject(s3, name)
	if err != nil {
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
	endpoint    string
	accessKey   string
	secretKey   string
	address     string
	bucket      string
	tlsCert     string
	tlsKey      string
	letsEncrypt bool
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "", "S3 server endpoint")
	flag.StringVar(&accessKey, "accessKey", "", "Access key of S3 storage")
	flag.StringVar(&secretKey, "secretKey", "", "Secret key of S3 storage")
	flag.StringVar(&bucket, "bucket", "", "Bucket name which hosts static files")
	flag.StringVar(&address, "address", "127.0.0.1:8080", "Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname")
	flag.StringVar(&tlsCert, "ssl-cert", "", "TLS certificate for this server")
	flag.StringVar(&tlsKey, "ssl-key", "", "TLS private key for this server")
	flag.BoolVar(&letsEncrypt, "lets-encrypt", false, "Enable Let's Encrypt")
}

// NewCustomHTTPTransport returns a new http configuration
// used while communicating with the cloud backends.
// This sets the value for MaxIdleConnsPerHost from 2 (go default)
// to 100.
func NewCustomHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          1024,
		MaxIdleConnsPerHost:   1024,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
	}
}

func main() {
	flag.Parse()

	if strings.TrimSpace(bucket) == "" {
		log.Fatalln(`Bucket name cannot be empty, please provide 's3www -bucket "mybucket"'`)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	// Chains all credential types, in the following order:
	//  - AWS env vars (i.e. AWS_ACCESS_KEY_ID)
	//  - AWS creds file (i.e. AWS_SHARED_CREDENTIALS_FILE or ~/.aws/credentials)
	//  - IAM profile based credentials. (performs an HTTP
	//    call to a pre-defined endpoint, only valid inside
	//    configured ec2 instances)
	var defaultAWSCredProviders = []credentials.Provider{
		&credentials.EnvAWS{},
		&credentials.FileAWSCredentials{},
		&credentials.IAM{
			Client: &http.Client{
				Transport: NewCustomHTTPTransport(),
			},
		},
		&credentials.EnvMinio{},
	}
	if accessKey != "" && secretKey != "" {
		defaultAWSCredProviders = []credentials.Provider{
			&credentials.Static{
				Value: credentials.Value{
					AccessKeyID:     accessKey,
					SecretAccessKey: secretKey,
				},
			},
		}
	}

	// If we see an Amazon S3 endpoint, then we use more ways to fetch backend credentials.
	// Specifically IAM style rotating credentials are only supported with AWS S3 endpoint.
	creds := credentials.NewChainCredentials(defaultAWSCredProviders)

	client, err := minio.NewWithOptions(u.Host, &minio.Options{
		Creds:        creds,
		Secure:       u.Scheme == "https",
		Region:       s3utils.GetRegionFromURL(*u),
		BucketLookup: minio.BucketLookupAuto,
	})
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.FileServer(&S3{client, bucket})
	if letsEncrypt {
		log.Printf("Started listening on https://%s\n", address)
		certmagic.HTTPS([]string{address}, mux)
	} else if tlsCert != "" && tlsKey != "" {
		log.Printf("Started listening on https://%s\n", address)
		log.Fatalln(http.ListenAndServeTLS(address, tlsCert, tlsKey, mux))
	} else {
		log.Printf("Started listening on http://%s\n", address)
		log.Fatalln(http.ListenAndServe(address, mux))
	}
}
