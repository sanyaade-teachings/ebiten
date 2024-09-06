// Copyright 2016 The Ebiten Authors
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

package restorable

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicscommand"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
)

// Image represents an image.
type Image struct {
	// Image is the underlying image.
	// This member is exported on purpose.
	// TODO: Move the implementation to internal/atlas package (#805).
	Image *graphicscommand.Image

	width  int
	height int
}

// NewImage creates an emtpy image with the given size.
//
// The returned image is cleared.
//
// Note that Dispose is not called automatically.
func NewImage(width, height int, screen bool) *Image {
	i := &Image{
		Image:  graphicscommand.NewImage(width, height, screen),
		width:  width,
		height: height,
	}

	// This needs to use 'InternalSize' to render the whole region, or edges are unexpectedly cleared on some
	// devices.
	iw, ih := i.Image.InternalSize()
	clearImage(i.Image, image.Rect(0, 0, iw, ih))
	return i
}

func clearImage(i *graphicscommand.Image, region image.Rectangle) {
	vs := make([]float32, 4*graphics.VertexFloatCount)
	graphics.QuadVerticesFromDstAndSrc(vs, float32(region.Min.X), float32(region.Min.Y), float32(region.Max.X), float32(region.Max.Y), 0, 0, 0, 0, 0, 0, 0, 0)
	is := graphics.QuadIndices()
	i.DrawTriangles([graphics.ShaderSrcImageCount]*graphicscommand.Image{}, vs, is, graphicsdriver.BlendClear, region, [graphics.ShaderSrcImageCount]image.Rectangle{}, clearShader.Shader, nil, graphicsdriver.FillRuleFillAll)
}

// ClearPixels clears the specified region by WritePixels.
func (i *Image) ClearPixels(region image.Rectangle) {
	if region.Dx() <= 0 || region.Dy() <= 0 {
		panic("restorable: width/height must be positive")
	}
	clearImage(i.Image, region.Intersect(image.Rect(0, 0, i.width, i.height)))
}
