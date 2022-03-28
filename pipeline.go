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
	"net/url"
	"strconv"
	"strings"
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

type Parser func(params []string, accept string) (Job, error)

func parseProxy(params []string, _ string) (Job, error) {
	if len(params) == 0 {
		return nil, errors.New("proxy url can not be nil")
	}
	proxyURL, err := url.QueryUnescape(params[0])
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, _ *Image) (*Image, error) {
		return FetchImageFromURL(ctx, proxyURL)
	}, nil
}

func parseOptimize(params []string, _ string) (Job, error) {
	if len(params) == 0 {
		return nil, errors.New("optimize addr can not be nil")
	}
	addr := params[0]
	quality := 0
	if len(params) > 1 {
		quality, _ = strconv.Atoi(params[1])
	}
	formats := make([]string, 0)

	if len(params) > 2 {
		formats = append(formats, params[2])
	}
	return NewOptimizeImage(addr, quality, formats...), nil
}

func parseAutoOptimize(params []string, accept string) (Job, error) {
	if len(params) == 0 {
		return nil, errors.New("optimize addr can not be nil")
	}
	quality := 0
	if len(params) > 1 {
		quality, _ = strconv.Atoi(params[1])
	}
	return NewAutoOptimizeImage(params[0], quality, accept), nil
}

func parseFitResize(params []string, _ string) (Job, error) {
	if len(params) != 2 {
		return nil, errors.New("fit resize width and height can not be nil")
	}
	// 如果转换出错，则直接用0
	width, _ := strconv.Atoi(params[0])
	height, _ := strconv.Atoi(params[1])
	return NewFitResizeImage(width, height), nil
}

func parseFillResize(params []string, _ string) (Job, error) {
	if len(params) != 2 {
		return nil, errors.New("fill resize width and height can not be nil")
	}
	// 如果转换出错，则直接用0
	width, _ := strconv.Atoi(params[0])
	height, _ := strconv.Atoi(params[1])
	return NewFillResizeImage(width, height), nil
}

func parseFinder(params []string, _ string) (Job, error) {
	if len(params) == 0 {
		return nil, errors.New("finder name can not be nil")
	}
	f, err := GetFinder(params[0])
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, _ *Image) (*Image, error) {
		return f.Find(ctx, params[1:]...)
	}, nil
}

const (
	TaskProxy        = "proxy"
	TaskOptimize     = "optimize"
	TaskAutoOptimize = "autoOptimize"
	TaskFitResize    = "fitResize"
	TaskFillResize   = "fillResize"
)

// Parse parses the task pipe line to job list
func Parse(taskPipeLine, accept string) ([]Job, error) {
	tasks := strings.Split(taskPipeLine, "|")
	jobs := make([]Job, 0, len(tasks))
	for _, v := range tasks {
		var fn Parser
		arr := strings.Split(v, "/")
		args := arr[1:]
		switch arr[0] {
		case TaskProxy:
			fn = parseProxy
		case TaskOptimize:
			fn = parseOptimize
		case TaskAutoOptimize:
			fn = parseAutoOptimize
		case TaskFitResize:
			fn = parseFitResize
		case TaskFillResize:
			fn = parseFillResize
		default:
			// finder的参数为所有参数
			args = arr
			fn = parseFinder
		}
		job, err := fn(args, accept)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
