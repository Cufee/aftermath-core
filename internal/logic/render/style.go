package render

import (
	"image/color"
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
	JustifyContentSpaceAround  // Spacing around each element is the same
	JustifyContentSpaceEvenly  // Similar to Grid, each element will have the same size
)

type directionValue int

const (
	DirectionHorizontal directionValue = iota
	DirectionVertical
)

type Style struct {
	JustifyContent justifyContentValue
	AlignItems     alignItemsValue // Depends on Direction
	Direction      directionValue

	Gap      float64
	PaddingX float64
	PaddingY float64

	Width  float64
	Height float64

	BorderRadius    float64
	BackgroundColor color.Color
}
