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

func TestNewFitResizeImage(t *testing.T) {
	assert := assert.New(t)

	width := 400
	height := 300
	fn := NewFitResizeImage(width, height)

	img, err := NewImageFromBytes(newImageData())
	assert.Nil(err)
	img, err = fn(context.Background(), img)
	assert.Nil(err)
	assert.Equal(293, img.Width())
	assert.Equal(height, img.Height())
}

func TestNewFillResizeImage(t *testing.T) {
	assert := assert.New(t)

	width := 400
	height := 300
	fn := NewFillResizeImage(width, height)

	img, err := NewImageFromBytes(newImageData())
	assert.Nil(err)
	img, err = fn(context.Background(), img)
	assert.Nil(err)
	assert.Equal(width, img.Width())
	assert.Equal(height, img.Height())
}
