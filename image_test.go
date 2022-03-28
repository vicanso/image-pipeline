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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageFromBytes(t *testing.T) {
	assert := assert.New(t)

	img, err := NewImageFromBytes(newImageData())
	assert.Nil(err)
	assert.Equal(829, img.Width())
	assert.Equal(846, img.Height())
}

func TestImageSet(t *testing.T) {
	assert := assert.New(t)

	img, err := NewImageFromBytes(newImageData())
	assert.Nil(err)

	img.Set(img.grid)
	assert.NotEqual(img.Previous(), img)
}

func TestImageSetOptimized(t *testing.T) {
	assert := assert.New(t)

	img, err := NewImageFromBytes(newImageData())
	assert.Nil(err)

	img.setOptimized([]byte("abc"), "jpeg")

	assert.Equal([]byte("abc"), img.optimizedData)
	assert.Equal("jpeg", img.format)
}

func TestImageEncode(t *testing.T) {
	assert := assert.New(t)

	img, err := NewImageFromBytes(newImageData())
	assert.Nil(err)

	buf, err := img.PNG()
	assert.Nil(err)
	assert.NotEmpty(buf)

	buf, err = img.JPEG()
	assert.Nil(err)
	assert.NotEmpty(buf)
}
