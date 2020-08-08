package main

import (
	"context"
	"os"
	"strings"
	"time"

	minio "github.com/minio/minio-go/v7"
)

const (
	pathSeparator = "/"
)

// A httpMinioObject implements http.File interface, returned by a S3
// Open method and can be served by the FileServer implementation.
type httpMinioObject struct {
	client *minio.Client
	object *minio.Object
	bucket string
	prefix string
	isDir  bool
}

func (h *httpMinioObject) Close() error {
	return h.object.Close()
}

func (h *httpMinioObject) Read(p []byte) (n int, err error) {
	return h.object.Read(p)
}

func (h *httpMinioObject) Seek(offset int64, whence int) (int64, error) {
	return h.object.Seek(offset, whence)
}

func (h *httpMinioObject) Readdir(count int) ([]os.FileInfo, error) {
	// List 'N' number of objects from a bucket-name with a matching prefix.
	listObjectsN := func(bucket, prefix string, count int) (objsInfo []minio.ObjectInfo, err error) {
		i := 1
		for object := range h.client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: false,
		}) {
			if object.Err != nil {
				return nil, object.Err
			}
			i++
			// Verify if we have printed N objects.
			if i == count {
				return
			}
			objsInfo = append(objsInfo, object)
		}
		return objsInfo, nil
	}

	// List non-recursively first count entries for prefix 'prefix" prefix.
	objsInfo, err := listObjectsN(h.bucket, h.prefix, count)
	if err != nil {
		return nil, os.ErrNotExist
	}
	var fileInfos []os.FileInfo
	for _, objInfo := range objsInfo {
		if strings.HasSuffix(objInfo.Key, pathSeparator) {
			fileInfos = append(fileInfos, objectInfo{
				ObjectInfo: minio.ObjectInfo{
					Key:          strings.TrimSuffix(objInfo.Key, pathSeparator),
					LastModified: objInfo.LastModified,
				},
				prefix: strings.TrimSuffix(objInfo.Key, pathSeparator),
				isDir:  true,
			})
			continue
		}
		fileInfos = append(fileInfos, objectInfo{
			ObjectInfo: objInfo,
		})
	}
	return fileInfos, nil
}

func (h *httpMinioObject) Stat() (os.FileInfo, error) {
	if h.isDir {
		return objectInfo{
			ObjectInfo: minio.ObjectInfo{
				Key:          h.prefix,
				LastModified: time.Now().UTC(),
			},
			prefix: h.prefix,
			isDir:  true,
		}, nil
	}

	objInfo, err := h.object.Stat()
	if err != nil {
		return nil, os.ErrNotExist
	}

	return objectInfo{
		ObjectInfo: objInfo,
	}, nil
}
