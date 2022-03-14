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
	"image"

	"github.com/disintegration/imaging"
)

type resizeHandler func(image.Image, int, int, imaging.ResampleFilter) *image.NRGBA

func resize(fn resizeHandler, img *Image, width, height int) (*Image, error) {
	w := img.Width()
	h := img.Height()
	if w <= width && h <= height {
		return img, nil
	}
	grid := fn(img.grid, width, height, imaging.Lanczos)
	img.Set(grid)
	return img, nil
}

// NewFitResizeImage creates an image job, which will resize the image to fit width/height
func NewFitResizeImage(width, height int) Job {
	return func(_ context.Context, img *Image) (*Image, error) {
		return resize(imaging.Fit, img, width, height)
	}
}

// NewFillResizeImage creates an image job, which will resize the image to fill with/height
func NewFillResizeImage(width, height int) Job {
	return func(_ context.Context, img *Image) (*Image, error) {
		return resize(func(i1 image.Image, i2, i3 int, rf imaging.ResampleFilter) *image.NRGBA {
			return imaging.Fill(i1, i2, i3, imaging.Center, rf)
		}, img, width, height)
	}
}
