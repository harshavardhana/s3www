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
	return os.FileMode(0o644)
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
