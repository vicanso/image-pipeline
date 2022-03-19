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
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/vicanso/tiny/pb"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcSF = singleflight.Group{}

var grpcConnections = sync.Map{}

var ErrGRPCClientInvalid = errors.New("grpc client connection is invalid")

func convertToConnection(value interface{}) (*grpc.ClientConn, error) {
	c, _ := value.(*grpc.ClientConn)
	if c == nil {
		return nil, ErrGRPCClientInvalid
	}
	return c, nil
}

func newTinyConnection(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	value, ok := grpcConnections.Load(addr)
	if ok {
		return convertToConnection(value)
	}
	value, err, _ := grpcSF.Do(addr, func() (interface{}, error) {
		conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		grpcConnections.Store(addr, conn)
		return conn, nil
	})
	if err != nil {
		return nil, err
	}
	return convertToConnection(value)
}

func optimize(ctx context.Context, addr string, img *Image, quality int, format string) (*Image, error) {
	c, err := newTinyConnection(ctx, addr)
	if err != nil {
		return nil, err
	}
	data, err := img.PNG()
	if err != nil {
		return nil, err
	}

	client := pb.NewOptimClient(c)
	in := pb.OptimRequest{
		Data:    data,
		Source:  pb.Type_PNG,
		Quality: uint32(quality),
	}
	switch format {
	case ImageTypePNG:
		in.Output = pb.Type_PNG
	case ImageTypeWEBP:
		in.Output = pb.Type_WEBP
	case ImageTypeAVIF:
		in.Output = pb.Type_AVIF
	default:
		in.Output = pb.Type_JPEG
	}
	reply, err := client.DoOptim(ctx, &in)
	if err != nil {
		return nil, err
	}
	img.setOptimized(reply.Data, format)
	return img, nil
}

// NewAutoOptimizeImage creates an optimize image job, which will find the match type for optimizing by accept
func NewAutoOptimizeImage(addr string, quality int, accept string) Job {
	return func(ctx context.Context, img *Image) (*Image, error) {
		format := img.format
		acceptWebp := strings.Contains(accept, "image/webp")
		acceptAvif := strings.Contains(accept, "image/avif")

		if acceptAvif {
			format = ImageTypeAVIF
		} else if acceptWebp {
			format = ImageTypeWEBP
			if format == ImageTypePNG {
				quality = 0
			}
		}
		return optimize(ctx, addr, img, quality, format)
	}
}

// NewOptimizeImage creates an optimize image job, it the format is nil, the original format will be used
func NewOptimizeImage(addr string, quality int, formats ...string) Job {
	return func(ctx context.Context, img *Image) (*Image, error) {
		format := img.format
		if len(formats) != 0 {
			format = formats[0]
		}
		return optimize(ctx, addr, img, quality, format)
	}
}
