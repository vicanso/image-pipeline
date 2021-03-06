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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWatermarkPosition(t *testing.T) {
	assert := assert.New(t)

	w := 800
	h := 600

	watermarkWidth := 60
	watermarkHeight := 40

	x, y := getWatermarkPosition(PositionTopLeft, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(0, x)
	assert.Equal(0, y)

	x, y = getWatermarkPosition(PositionTop, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(370, x)
	assert.Equal(0, y)

	x, y = getWatermarkPosition(PositionTopRight, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(740, x)
	assert.Equal(0, y)

	x, y = getWatermarkPosition(PositionLeft, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(0, x)
	assert.Equal(280, y)

	x, y = getWatermarkPosition(PositionCenter, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(370, x)
	assert.Equal(280, y)

	x, y = getWatermarkPosition(PositionRight, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(740, x)
	assert.Equal(280, y)

	x, y = getWatermarkPosition(PositionBottomLeft, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(0, x)
	assert.Equal(560, y)

	x, y = getWatermarkPosition(PositionBottom, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(370, x)
	assert.Equal(560, y)

	x, y = getWatermarkPosition(PositionBottomRight, w, h, watermarkWidth, watermarkHeight)
	assert.Equal(740, x)
	assert.Equal(560, y)
}

func TestNewWatermark(t *testing.T) {
	assert := assert.New(t)

	img, err := NewImageFromBytes(newImageData())
	assert.Nil(err)

	watermark, err := NewFitResizeImage(100, 100)(context.Background(), img)
	assert.Nil(err)

	img, err = NewImageFromBytes(newImageData())
	assert.Nil(err)
	fn := NewWatermark(watermark.grid, PositionBottom, 0)
	img, err = fn(context.Background(), img)
	assert.Nil(err)
	assert.Equal(829, img.Width())
	assert.Equal(846, img.Height())
}
