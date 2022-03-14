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
)

const (
	PositionTopLeft     = "topLeft"
	PositionTop         = "top"
	PositionTopRight    = "topRight"
	PositionLeft        = "left"
	PositionCenter      = "center"
	PositionRight       = "right"
	PositionBottomLeft  = "bottomLeft"
	PositionBottom      = "bottom"
	PositionBottomRight = "bottomRight"
)

const (
	ImageTypePNG  = "png"
	ImageTypeJPEG = "jpeg"
	ImageTypeWEBP = "webp"
	ImageTypeAVIF = "avif"
)

type ImageFinder func(ctx context.Context, params ...string) (*Image, error)

// 不再执行后续时返回
var ErrAbortNext = errors.New("abort next")

// Do runs the pipeline jobs
func Do(ctx context.Context, img *Image, jobs ...Job) (*Image, error) {
	var err error
	for _, fn := range jobs {
		img, err = fn(ctx, img)
		if err != nil {
			// 如果是abort error，则直接返回数据
			if err == ErrAbortNext {
				return img, nil
			}
			return nil, err
		}
	}
	return img, nil
}
