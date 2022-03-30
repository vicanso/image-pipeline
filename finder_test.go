// Copyright 2022 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package imagepipeline

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAddHTTPFinder(t *testing.T) {

	assert := assert.New(t)
	ln, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer ln.Close()

	mux := http.NewServeMux()
	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}))

	mux.Handle("/image", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := newImageData()
		w.Header().Add("Content-Type", "mage/png")
		_, _ = w.Write(buf)
	}))

	s := http.Server{
		Handler: mux,
	}
	go func() {
		_ = s.Serve(ln)
	}()
	time.Sleep(100 * time.Millisecond)
	uri := fmt.Sprintf("http://%s/ping", ln.Addr().String())
	finderName := "httpFinder"
	err = AddHTTPFinder(finderName, uri)
	assert.Nil(err)
	finder, err := GetFinder(finderName)
	assert.Nil(err)
	assert.NotNil(finder)

	img, err := finder.Find(context.Background(), "/image")
	assert.Nil(err)
	assert.Equal(829, img.Width())
	assert.Equal(846, img.Height())

	err = finder.Close(context.Background())
	assert.Nil(err)
}

func TestFileFinder(t *testing.T) {
	assert := assert.New(t)
	basePath := os.TempDir()
	finderName := "fileFinder"
	err := AddFileFinder(finderName, basePath)
	assert.Nil(err)

	file := "/test.png"
	err = ioutil.WriteFile(basePath+file, newImageData(), 0777)
	assert.Nil(err)

	finder, err := GetFinder(finderName)
	assert.Nil(err)
	assert.NotNil(finder)

	img, err := finder.Find(context.Background(), file)
	assert.Nil(err)
	assert.Equal(829, img.Width())
	assert.Equal(846, img.Height())
}

func TestMinioFinder(t *testing.T) {
	assert := assert.New(t)
	finderName := "minioFinder"
	host := "127.0.0.1:9000"
	accessKey := "test"
	secretKey := "testabcd"
	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	assert.Nil(err)
	bucket := "findertest"
	objName := "test.png"
	exist, err := client.BucketExists(context.Background(), bucket)
	assert.Nil(err)

	if !exist {
		err = client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
		assert.Nil(err)
		buf := newImageData()
		_, err = client.PutObject(context.Background(), bucket, objName, bytes.NewReader(buf), int64(len(buf)), minio.PutObjectOptions{})
		assert.Nil(err)
	}

	err = AddMinioFinder(finderName, fmt.Sprintf("http://%s/?accessKey=%s&secretKey=%s", host, accessKey, secretKey))
	assert.Nil(err)

	finder, err := GetFinder(finderName)
	assert.Nil(err)
	assert.NotNil(finder)

	img, err := finder.Find(context.Background(), bucket, objName)
	assert.Nil(err)
	assert.Equal(829, img.Width())
	assert.Equal(846, img.Height())
}

func TestGridFSFinder(t *testing.T) {
	assert := assert.New(t)

	finderName := "gridFinder"
	database := "admin"
	uri := "mongodb://test:testabcd@127.0.0.1:27017/" + database
	fileName := "test.png"

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	assert.Nil(err)
	bucket, err := gridfs.NewBucket(client.Database(database), options.GridFSBucket().SetName(options.DefaultName))
	assert.Nil(err)
	stream, err := bucket.OpenUploadStream(fileName)
	assert.Nil(err)
	_, err = stream.Write(newImageData())
	assert.Nil(err)
	err = stream.Close()
	assert.Nil(err)

	err = AddGridFSFinder(finderName, uri)
	assert.Nil(err)

	finder, err := GetFinder(finderName)
	assert.Nil(err)
	assert.NotNil(finder)
	defer finder.Close(context.Background())
	id, ok := stream.FileID.(primitive.ObjectID)
	assert.True(ok)
	img, err := finder.Find(context.Background(), id.Hex())
	assert.Nil(err)
	assert.Equal(829, img.Width())
	assert.Equal(846, img.Height())
}

func TestAliyunOSSFinder(t *testing.T) {
	assert := assert.New(t)

	accessKey := os.Getenv("ALIYUN_OSS_ACCESS_KEY")
	secretKey := os.Getenv("ALIYUN_OSS_SECRET_KEY")
	uri := fmt.Sprintf("https://oss-cn-beijing.aliyuncs.com?accessKey=%s&secretKey=%s", accessKey, secretKey)

	finderName := "aliyunOSSFinder"
	err := AddAliyunOSSFinder(finderName, uri)
	assert.Nil(err)
	finder, err := GetFinder(finderName)
	assert.Nil(err)
	assert.NotNil(finder)

	img, err := finder.Find(context.Background(), "tinysite", "go-echarts.jpg")
	assert.Nil(err)
	assert.Equal(829, img.Width())
	assert.Equal(846, img.Height())
}
