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
	"errors"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vicanso/upstream"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var ErrFinderNotFound = errors.New("Finder is not found")
var ErrFinderInValid = errors.New("Finder is invald")

type Finder interface {
	Find(ctx context.Context, params ...string) (*Image, error)
	Close(ctx context.Context) error
}

type finder struct {
	Find  func(ctx context.Context, params ...string) (*Image, error)
	Close func(ctx context.Context) error
}

func noop(ctx context.Context) error {
	return nil
}

var finders = sync.Map{}

// AddHTTPFinder adds a http finder
func AddHTTPFinder(name, uri string, onStatus ...upstream.StatusListener) error {
	urlInfo, err := url.Parse(uri)
	if err != nil {
		return err
	}
	uh := &upstream.HTTP{
		Ping: urlInfo.Path,
	}
	for _, host := range strings.Split(urlInfo.Host, ",") {
		uh.Add(urlInfo.Scheme + "://" + host)
	}
	if len(onStatus) != 0 {
		uh.OnStatus(onStatus[0])
	}

	go uh.StartHealthCheck()
	f := finder{
		Find: func(ctx context.Context, params ...string) (*Image, error) {
			if len(params) < 1 {
				return nil, errors.New("http params should be on parameter")
			}
			requestURI, err := url.QueryUnescape(params[0])
			if err != nil {
				return nil, err
			}
			u := uh.PolicyRoundRobin()
			if u == nil {
				return nil, errors.New("get http upstream fail")
			}
			return FetchImageFromURL(ctx, u.URL.String()+requestURI)
		},
		Close: func(ctx context.Context) error {
			uh.StopHealthCheck()
			return nil
		},
	}
	finders.Store(name, &f)
	return nil
}

// AddMinioFinder adds a minio finder
func AddMinioFinder(name, uri string) error {
	urlInfo, err := url.Parse(uri)
	if err != nil {
		return err
	}
	accessKey := urlInfo.Query().Get("accessKey")
	secretKey := urlInfo.Query().Get("secretKey")
	minioClient, err := minio.New(urlInfo.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	f := finder{
		Find: func(ctx context.Context, params ...string) (*Image, error) {
			if len(params) < 2 {
				return nil, errors.New("minio param should be two parameters")
			}
			obj, err := minioClient.GetObject(ctx, params[0], params[1], minio.GetObjectOptions{})
			if err != nil {
				return nil, err
			}
			buf, err := io.ReadAll(obj)
			if err != nil {
				return nil, err
			}
			return NewImageFromBytes(buf)
		},
		Close: noop,
	}
	finders.Store(name, &f)
	return nil
}

// AddGridFSFinder adds mongodb gridfs finder
func AddGridFSFinder(name, uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return err
	}
	if cs.Database == "" {
		return errors.New("database can not be nil")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	f := finder{
		Find: func(ctx context.Context, params ...string) (*Image, error) {
			if len(params) == 0 {
				return nil, errors.New("gridfs param should be one parameter")
			}
			db := client.Database(cs.Database)
			collection := options.DefaultName
			if len(params) > 1 {
				collection = params[1]
			}
			bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(collection))
			if err != nil {
				return nil, err
			}
			buffer := bytes.Buffer{}
			id, err := primitive.ObjectIDFromHex(params[0])
			if err != nil {
				return nil, err
			}
			_, err = bucket.DownloadToStream(id, &buffer)
			if err != nil {
				return nil, err
			}
			return NewImageFromBytes(buffer.Bytes())
		},
		Close: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	}
	finders.Store(name, &f)
	return nil
}

// AddAliyunOSSFinder add aliyun oss finder
func AddAliyunOSSFinder(name, uri string) error {
	urlInfo, err := url.Parse(uri)
	if err != nil {
		return err
	}
	accessKey := urlInfo.Query().Get("accessKey")
	secretKey := urlInfo.Query().Get("secretKey")
	if len(accessKey) == 0 || len(secretKey) == 0 {
		return errors.New("access key and secret key can not be nil")
	}
	client, err := oss.New(urlInfo.Hostname(), accessKey, secretKey)
	if err != nil {
		return err
	}
	f := finder{
		Find: func(_ context.Context, params ...string) (*Image, error) {
			if len(params) < 2 {
				return nil, errors.New("oss param should be two parameters")
			}
			bucket, err := client.Bucket(params[0])
			if err != nil {
				return nil, err
			}
			r, err := bucket.GetObject(params[1])
			if err != nil {
				return nil, err
			}
			buf, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}

			return NewImageFromBytes(buf)
		},
		Close: noop,
	}
	finders.Store(name, &f)
	return nil
}

// GetFinder returns a finder by name
func GetFinder(name string) (Finder, error) {
	value, ok := finders.Load(name)
	if !ok {
		return nil, ErrFinderNotFound
	}
	fn, ok := value.(Finder)
	if !ok {
		return nil, ErrFinderInValid
	}
	return fn, nil
}

// RangeFinder calls fn sequentially for each key and find in the finders
func RangeFinder(fn func(name string, f Finder)) {
	finders.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		f, ok := value.(Finder)
		if ok {
			fn(k, f)
		}
		return true
	})
}
