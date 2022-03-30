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
	"image/color"

	"github.com/disintegration/imaging"
)

func getWatermarkPosition(position string, w, h, watermarkWidth, watermarkHeight int) (int, int) {
	x := 0
	y := 0
	// PositionTopLeft 为0,0 不需要处理
	switch position {
	case PositionTop:
		x = (w - watermarkWidth) / 2
	case PositionTopRight:
		x = w - watermarkWidth
	case PositionLeft:
		y = (h - watermarkHeight) / 2
	case PositionCenter:
		x = (w - watermarkWidth) / 2
		y = (h - watermarkHeight) / 2
	case PositionRight:
		y = (h - watermarkHeight) / 2
		x = w - watermarkWidth
	case PositionBottomLeft:
		y = h - watermarkHeight
	case PositionBottom:
		x = (w - watermarkWidth) / 2
		y = h - watermarkHeight
	case PositionBottomRight:
		x = w - watermarkWidth
		y = h - watermarkHeight
	}
	return x, y
}

// NewWatermark creates an image job, which will add watermark to image
func NewWatermark(watermarkImg image.Image, position string, angle float64) Job {
	return func(ctx context.Context, img *Image) (*Image, error) {
		if angle != 0 {
			watermarkImg = imaging.Rotate(watermarkImg, angle, color.Transparent)
		}
		w := img.Width()
		h := img.Height()
		watermarkWidth := watermarkImg.Bounds().Dx()
		watermarkHeight := watermarkImg.Bounds().Dy()
		x, y := getWatermarkPosition(position, w, h, watermarkWidth, watermarkHeight)
		grid := imaging.Paste(img.grid, watermarkImg, image.Pt(x, y))
		img.Set(grid)
		return img, nil
	}
}
