package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	minio "github.com/minio/minio-go"
)

type s3FS struct {
	mc         *minio.Client
	bucketName string
}

type httpFile struct {
	client     *minio.Client
	object     *minio.Object
	bucketName string
	prefixName string
	isDir      bool
}

type mObjectInfo struct {
	objInfo    minio.ObjectInfo
	prefixName string
	isDir      bool
}

func (o mObjectInfo) Name() string {
	return o.objInfo.Key
}

func (o mObjectInfo) Size() int64 {
	return o.objInfo.Size
}

func (o mObjectInfo) Mode() os.FileMode {
	if o.isDir {
		return os.ModeDir
	}
	return os.FileMode(0644)
}

func (o mObjectInfo) ModTime() time.Time {
	return o.objInfo.LastModified
}

func (o mObjectInfo) IsDir() bool {
	return o.isDir
}

func (o mObjectInfo) Sys() interface{} {
	return &syscall.Stat_t{}
}

func (h *httpFile) Close() error {
	return h.object.Close()
}

func (h *httpFile) Read(p []byte) (n int, err error) {
	return h.object.Read(p)
}

func (h *httpFile) Seek(offset int64, whence int) (int64, error) {
	return h.object.Seek(offset, whence)
}

func (h *httpFile) Readdir(count int) ([]os.FileInfo, error) {
	if h.bucketName == "/" {
		buckets, err := h.client.ListBuckets()
		if err != nil {
			return nil, err
		}
		var fileInfos []os.FileInfo
		for _, bucket := range buckets {
			fileInfos = append(fileInfos, mObjectInfo{
				isDir: true,
				objInfo: minio.ObjectInfo{
					Key:          bucket.Name,
					LastModified: bucket.CreationDate,
				},
				prefixName: bucket.Name,
			})
		}
		return fileInfos, nil
	}
	// List 'N' number of objects from a bucket-name with a matching prefix.
	listObjectsN := func(bucket, prefix string, recursive bool, N int) (objsInfo []minio.ObjectInfo, err error) {
		// Create a done channel to control 'ListObjects' go routine.
		doneCh := make(chan struct{}, 1)

		// Free the channel upon return.
		defer close(doneCh)

		i := 1
		for object := range h.client.ListObjects(bucket, prefix, recursive, doneCh) {
			if object.Err != nil {
				return nil, object.Err
			}
			i++
			// Verify if we have printed N objects.
			if i == N {
				// Indicate ListObjects go-routine to exit and stop
				// feeding the objectInfo channel.
				doneCh <- struct{}{}
			}
			objsInfo = append(objsInfo, object)
		}
		return objsInfo, nil
	}
	// List non-recursively first count entries for prefix 'prefixName" prefix.
	recursive := count == -1
	objsInfo, err := listObjectsN(h.bucketName, h.prefixName, recursive, count)
	if err != nil {
		return nil, err
	}
	var fileInfos []os.FileInfo
	for _, objInfo := range objsInfo {
		if strings.HasSuffix(objInfo.Key, "/") {
			fileInfos = append(fileInfos, mObjectInfo{
				isDir: true,
				objInfo: minio.ObjectInfo{
					Key:          objInfo.Key,
					LastModified: time.Now().UTC(),
				},
				prefixName: objInfo.Key,
			})
			continue
		}
		fileInfos = append(fileInfos, mObjectInfo{objInfo: objInfo})
	}
	return fileInfos, nil
}

func (h *httpFile) Stat() (os.FileInfo, error) {
	if h.isDir {
		return mObjectInfo{isDir: h.isDir, prefixName: h.prefixName}, nil
	}
	objInfo, err := h.object.Stat()
	if err != nil {
		return nil, err
	}
	return mObjectInfo{objInfo: objInfo}, nil
}

func newHTTPFile(client *minio.Client, mObject *minio.Object, bucketName string, prefixName string) http.File {
	return &httpFile{
		object:     mObject,
		client:     client,
		isDir:      mObject == nil,
		bucketName: bucketName,
		prefixName: prefixName,
	}
}

func (s3 *s3FS) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, "/") {
		return newHTTPFile(s3.mc, nil, s3.bucketName, path.Clean(name)), nil
	}
	name = strings.TrimPrefix(name, "/")
	obj, err := s3.mc.GetObject(s3.bucketName, name)
	if err != nil {
		return nil, err
	}
	return newHTTPFile(s3.mc, obj, s3.bucketName, name), nil
}

func newS3FS(mc *minio.Client, bucketName string) *s3FS {
	return &s3FS{
		mc:         mc,
		bucketName: bucketName,
	}
}

// Convert string to bool and always return true if any error
func mustParseBool(str string) bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return true
	}
	return b
}

var (
	endpoint  string
	accessKey string
	secretKey string
	address   string
	bucket    string
	secure    bool
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "play.minio.io:9000", "S3 compatible storage endpoint")
	flag.StringVar(&accessKey, "accessKey", "Q3AM3UQ867SPQQA43P2F", "Access key")
	flag.StringVar(&secretKey, "secretKey", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "Secret key")
	flag.StringVar(&address, "address", ":8080", "Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname.")
	flag.StringVar(&bucket, "bucket", "", "Bucket name which hosts static files")
	flag.BoolVar(&secure, "s", false, "Enables secure communication with S3 compatible storage")
}

func main() {
	flag.Parse()

	mc, err := minio.New(endpoint, accessKey, secretKey, secure)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Started listening to ", address)
	log.Fatalln(http.ListenAndServe(address, http.FileServer(newS3FS(mc, bucket))))
}
