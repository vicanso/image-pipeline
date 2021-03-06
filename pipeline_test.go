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

func TestParseProxy(t *testing.T) {
	assert := assert.New(t)

	_, err := parseProxy([]string{
		"https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png",
	}, "")
	assert.Nil(err)
}

func TestParseOptimize(t *testing.T) {
	assert := assert.New(t)

	_, err := parseOptimize([]string{
		"127.0.0.1:6002",
		"90",
		"png",
	}, "")
	assert.Nil(err)
}

func TestParseAutoOptimize(t *testing.T) {
	assert := assert.New(t)

	_, err := parseAutoOptimize([]string{
		"127.0.0.1:6002",
		"80",
	}, "image/avif,image/webp")
	assert.Nil(err)
}

func TestParseFitResize(t *testing.T) {
	assert := assert.New(t)

	_, err := parseFitResize([]string{
		"100",
		"80",
	}, "")
	assert.Nil(err)
}

func TestParseFillResize(t *testing.T) {
	assert := assert.New(t)

	_, err := parseFillResize([]string{
		"100",
		"80",
	}, "")
	assert.Nil(err)
}

func TestParseWatermark(t *testing.T) {
	assert := assert.New(t)

	_, err := parseWatermark([]string{
		"https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png",
		"bottomRight",
	}, "")
	assert.Nil(err)
}
