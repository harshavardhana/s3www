// Copyright 2021 Harshavardhana
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/caddyserver/certmagic"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/rs/cors"
)

// Use e.g.: go build -ldflags "-X main.version=v1.0.0"
// to set the binary version.
var version = "0.0.0-dev"

// S3 - A S3 implements FileSystem using the minio client
// allowing access to your S3 buckets and objects.
//
// Note that S3 will allow all access to files in your private
// buckets, If you have any sensitive information please make
// sure to not sure this project.
type S3 struct {
	*minio.Client
	bucket     string
	bucketPath string
}

func pathIsDir(ctx context.Context, s3 *S3, name string) bool {
	name = strings.Trim(name, pathSeparator) + pathSeparator
	listCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	objCh := s3.Client.ListObjects(listCtx,
		s3.bucket,
		minio.ListObjectsOptions{
			Prefix:  name,
			MaxKeys: 1,
		})
	for range objCh {
		cancel()
		return true
	}
	return false
}

// Open - implements http.Filesystem implementation.
func (s3 *S3) Open(name string) (http.File, error) {
	name = path.Join(s3.bucketPath, name)
	if name == pathSeparator || pathIsDir(context.Background(), s3, name) {
		return &httpMinioObject{
			client: s3.Client,
			object: nil,
			isDir:  true,
			bucket: bucket,
			prefix: strings.TrimSuffix(name, pathSeparator),
		}, nil
	}

	name = strings.TrimPrefix(name, pathSeparator)
	obj, err := getObject(context.Background(), s3, name)
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

func getObject(ctx context.Context, s3 *S3, name string) (*minio.Object, error) {
	names := []string{name, name + "/index.html", name + "/index.htm"}
	if spaFile != "" {
		names = append(names, spaFile)
	}
	names = append(names, "/404.html")
	for _, n := range names {
		obj, err := s3.Client.GetObject(ctx, s3.bucket, n, minio.GetObjectOptions{})
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = obj.Stat()
		if err != nil {
			// do not log "file" in bucket not found errors
			if minio.ToErrorResponse(err).Code != "NoSuchKey" {
				log.Println(err)
			}
			continue
		}

		return obj, nil
	}

	return nil, os.ErrNotExist
}

var (
	endpoint          string
	accessKey         string
	accessKeyFile     string
	secretKey         string
	secretKeyFile     string
	address           string
	bucket            string
	bucketPath        string
	tlsCert           string
	tlsKey            string
	spaFile           string
	allowedCorsOrigin string
	letsEncrypt       bool
	versionF          = flag.Bool("version", false, "print version")
)

func init() {
	flag.BoolVar(versionF, "v", false, "print version")
	flag.StringVar(&endpoint, "endpoint", defaultEnvString("S3WWW_ENDPOINT", ""), "AWS S3 compatible server endpoint")
	flag.StringVar(&bucket, "bucket", defaultEnvString("S3WWW_BUCKET", ""), "bucket name with static files")
	flag.StringVar(&bucketPath, "bucketPath", defaultEnvString("S3WWW_BUCKET_PATH", "/"), "bucket path to serve static files from")
	flag.StringVar(&accessKey, "accessKey", defaultEnvString("S3WWW_ACCESS_KEY", ""), "access key for server")
	flag.StringVar(&secretKey, "secretKey", defaultEnvString("S3WWW_SECRET_KEY", ""), "secret key for server")
	flag.StringVar(&address, "address", defaultEnvString("S3WWW_ADDRESS", "127.0.0.1:8080"), "bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname")
	flag.BoolVar(&letsEncrypt, "lets-encrypt", defaultEnvBool("S3WWW_LETS_ENCRYPT", false), "enable Let's Encrypt for automatic TLS certs for the DOMAIN")
	flag.StringVar(&tlsCert, "ssl-cert", defaultEnvString("S3WWW_SSL_CERT", ""), "public TLS certificate for this server")
	flag.StringVar(&tlsKey, "ssl-key", defaultEnvString("S3WWW_SSL_KEY", ""), "private TLS key for this server")
	flag.StringVar(&accessKeyFile, "accessKeyFile", defaultEnvString("S3WWW_ACCESS_KEY_FILE", ""), "file which contains the access key")
	flag.StringVar(&secretKeyFile, "secretKeyFile", defaultEnvString("S3WWW_SECRET_KEY_FILE", ""), "file which contains the secret key")
	flag.StringVar(&spaFile, "spaFile", defaultEnvString("S3WWW_SPA_FILE", ""), "if working with SPA (Single Page Application), use this key the set the absolute path of the file to call whenever a file dosen't exist")
	flag.StringVar(&allowedCorsOrigin, "allowed-cors-origins", defaultEnvString("S3WWW_ALLOWED_CORS_ORIGINS", ""), "a list of origins a cross-domain request can be executed from")
}

func defaultEnvString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func defaultEnvBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		parsedVal, err := strconv.ParseBool(val)
		if err == nil {
			return parsedVal
		}
		log.Printf("String of %q did not parse as bool for env var %q", val, key)
	}
	return defaultVal
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
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		DisableCompression:    true,
	}
}

func main() {
	flag.Parse()

	if *versionF {
		fmt.Println("s3www -", version)
		os.Exit(0)
	}

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
	defaultAWSCredProviders := []credentials.Provider{
		&credentials.EnvAWS{},
		&credentials.FileAWSCredentials{},
		&credentials.IAM{
			Client: &http.Client{
				Transport: NewCustomHTTPTransport(),
			},
		},
		&credentials.EnvMinio{},
	}
	if accessKeyFile != "" {
		if keyBytes, err := os.ReadFile(accessKeyFile); err == nil {
			accessKey = strings.TrimSpace(string(keyBytes))
		} else {
			log.Fatalf("Failed to read access key file %q", accessKeyFile)
		}
	}
	if secretKeyFile != "" {
		if keyBytes, err := os.ReadFile(secretKeyFile); err == nil {
			secretKey = strings.TrimSpace(string(keyBytes))
		} else {
			log.Fatalf("Failed to read secret key file %q", secretKeyFile)
		}
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

	client, err := minio.New(u.Host, &minio.Options{
		Creds:        creds,
		Secure:       u.Scheme == "https",
		Region:       s3utils.GetRegionFromURL(*u),
		BucketLookup: minio.BucketLookupAuto,
		Transport:    NewCustomHTTPTransport(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.FileServer(&S3{client, bucket, bucketPath})

	// Wrap the existing mux with the CORS middleware.
	opts := cors.Options{
		AllowOriginFunc: func(origin string) bool {
			if allowedCorsOrigin == "" {
				return true
			}
			for _, allowedOrigin := range strings.Split(allowedCorsOrigin, ",") {
				if wildcard.Match(allowedOrigin, origin) {
					return true
				}
			}
			return false
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodHead,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodPatch,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
	muxHandler := cors.New(opts).Handler(mux)

	switch {
	case letsEncrypt:
		log.Printf("Started listening on https://%s\n", address)
		certmagic.HTTPS([]string{address}, muxHandler)
	case tlsCert != "" && tlsKey != "":
		log.Printf("Started listening on https://%s\n", address)
		log.Fatalln(http.ListenAndServeTLS(address, tlsCert, tlsKey, muxHandler))
	default:
		log.Printf("Started listening on http://%s\n", address)
		log.Fatalln(http.ListenAndServe(address, muxHandler))
	}
}
