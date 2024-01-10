package render

type alignItemsValue int

const (
	AlignItemsStart alignItemsValue = iota
	AlignItemsCenter
	AlignItemsEnd
)

type directionValue int

const (
	DirectionHorizontal directionValue = iota
	DirectionVertical
)

type Style struct {
	AlignItems alignItemsValue // Depends on Direction
	Direction  directionValue

	Gap     float64
	Padding float64 // not implemented
}
