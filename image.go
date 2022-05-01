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
	"image"

	"github.com/disintegration/imaging"
)

type Image struct {
	// previous is the previous before handle
	previous *Image
	// originalSize is the original raw data size of image
	originalSize int
	// grid is the grid of color.Color values
	grid image.Image
	// optimizedData is the data of optimize image
	optimizedData []byte
	// format is the format type of image
	format string
}

// Job is the image pipeline job
type Job func(context.Context, *Image) (*Image, error)

// NewImageFromBytes returns a image from byte data, an error will be return if decode fail
func NewImageFromBytes(data []byte) (*Image, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &Image{
		optimizedData: data,
		originalSize:  len(data),
		format:        format,
		grid:          img,
	}, nil
}

// Previous returns the previous image
func (i *Image) Previous() *Image {
	return i.previous
}

// Set sets the image grid
func (i *Image) Set(grid image.Image) {
	previous := *i
	i.previous = &previous
	// the image is changed, reset the optimized data
	i.optimizedData = nil
	i.grid = grid
}

// Width returns the width of image
func (i *Image) Width() int {
	return i.grid.Bounds().Dx()
}

// Height returns the height of image
func (i *Image) Height() int {
	return i.grid.Bounds().Dy()
}

func (i *Image) setOptimized(data []byte, format string) {
	i.optimizedData = data
	i.format = format
}

func (i *Image) encode(format string) ([]byte, error) {
	buffer := bytes.Buffer{}
	f := imaging.JPEG
	if format == ImageTypePNG {
		f = imaging.PNG
	}
	err := imaging.Encode(&buffer, i.grid, f)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// PNG encodes the image as png, and returns the bytes
func (i *Image) PNG() ([]byte, error) {
	if i.format == ImageTypePNG &&
		len(i.optimizedData) != 0 {
		return i.optimizedData, nil
	}
	return i.encode(ImageTypePNG)
}

// JPEG encodes the image as jpeg, and returns the bytes
func (i *Image) JPEG() ([]byte, error) {
	if i.format == ImageTypeJPEG &&
		len(i.optimizedData) != 0 {
		return i.optimizedData, nil
	}
	return i.encode(ImageTypeJPEG)
}

// Bytes returns the bytes and format of image
func (i *Image) Bytes() ([]byte, string) {
	return i.optimizedData, i.format
}
