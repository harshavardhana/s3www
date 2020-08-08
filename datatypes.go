package main

import (
	"os"
	"syscall"
	"time"

	minio "github.com/minio/minio-go/v7"
)

// objectInfo implements os.FileInfo interface,
// is returned during Readdir(), Stat() operations.
type objectInfo struct {
	minio.ObjectInfo
	prefix string
	isDir  bool
}

func (o objectInfo) Name() string {
	return o.ObjectInfo.Key
}

func (o objectInfo) Size() int64 {
	return o.ObjectInfo.Size
}

func (o objectInfo) Mode() os.FileMode {
	if o.isDir {
		return os.ModeDir
	}
	return os.FileMode(0644)
}

func (o objectInfo) ModTime() time.Time {
	return o.ObjectInfo.LastModified
}

func (o objectInfo) IsDir() bool {
	return o.isDir
}

func (o objectInfo) Sys() interface{} {
	return &syscall.Stat_t{}
}
