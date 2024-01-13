package render

import (
	"image/color"

	"golang.org/x/image/font"
)

type alignItemsValue int
type justifyContentValue int

const (
	AlignItemsStart alignItemsValue = iota
	AlignItemsCenter
	AlignItemsEnd

	JustifyContentStart justifyContentValue = iota
	JustifyContentCenter
	JustifyContentEnd
	JustifyContentSpaceBetween // Spacing between each element is the same
)

type directionValue int

const (
	DirectionHorizontal directionValue = iota
	DirectionVertical
)

type Style struct {
	Font      font.Face
	FontColor color.Color

	JustifyContent justifyContentValue
	AlignItems     alignItemsValue // Depends on Direction
	Direction      directionValue

	Gap float64

	PaddingX float64
	PaddingY float64

	Width  float64
	Height float64

	BorderRadius    float64
	BackgroundColor color.Color

	Debug bool
}
