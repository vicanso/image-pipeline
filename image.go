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
	"image"
)

type Image struct {
	// previous is the previous before handle
	previous *Image
	// data is the raw data of image
	data []byte
	img  image.Image
}

// NewImageFromBytes returns a image from byte data, an error will be return if decode fail
func NewImageFromBytes(data []byte) (*Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &Image{
		data: data,
		img:  img,
	}, nil
}

// Previous returns the previous image
func (i *Image) Previous() *Image {
	return i.previous
}

// Image returns image interface
func (i *Image) Image() (image.Image, error) {
	if i.img == nil {
		img, _, err := image.Decode(bytes.NewReader(i.data))
		if err != nil {
			return nil, err
		}
		i.img = img
	}
	return i.img, nil
}

// Data returns the raw data of image
func (i *Image) Data() []byte {
	return i.data
}

// Size returns the size of image's data
func (i *Image) Size() int {
	return len(i.data)
}

func (i *Image) bounds() (*image.Rectangle, error) {
	img, err := i.Image()
	if err != nil {
		return nil, err
	}
	r := img.Bounds()
	return &r, nil
}

// Width returns the width of image
func (i *Image) Width() (int, error) {
	b, err := i.bounds()
	if err != nil {
		return 0, err
	}
	return b.Dx(), nil
}

// Height returns the height of image
func (i *Image) Height() (int, error) {
	b, err := i.bounds()
	if err != nil {
		return 0, err
	}
	return b.Dy(), nil
}

// SetData sets the data of image, and clear the image interface
func (i *Image) SetData(data []byte) {
	i.data = data
	i.img = nil
}
