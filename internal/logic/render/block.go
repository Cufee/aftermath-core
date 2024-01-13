package render

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/fogleman/gg"
)

var (
	debugColorPink = color.RGBA{255, 192, 203, 255}
)

type block struct {
	content []blockContent
	options *RenderOptions
}

func (block *block) Render() (image.Image, error) {
	var images []image.Image
	for _, row := range block.content {
		img, err := row.Render()
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return renderImages(images, block.options)
}

type blockContent struct {
	value   string
	options *RenderOptions
}

func (row *blockContent) Render() (image.Image, error) {
	if row.options.Style.Font == nil {
		return nil, errors.New("font not set")
	}

	measureCtx := gg.NewContext(1, 1)
	measureCtx.SetFontFace(row.options.Style.Font)
	valueW, valueH := measureCtx.MeasureString(row.value)

	// Account for font descender height
	descenderOffset := (float64(row.options.Style.Font.Metrics().Descent>>6) - 1)
	ctx := gg.NewContext(int(math.Ceil(valueW)+1), int(math.Ceil(valueH+(descenderOffset*2))))

	// Render text
	ctx.SetFontFace(row.options.Style.Font)
	ctx.SetColor(row.options.Style.FontColor)

	ctx.DrawString(row.value, 0, valueH)

	return ctx.Image(), nil
}

func NewBlock(label string, options *RenderOptions, rows ...any) block {
	var contentRows []blockContent
	for _, row := range rows {
		value := validOrPlaceholder(row)
		if value == "" {
			continue
		}
		contentRows = append(contentRows, blockContent{
			value:   value,
			options: options,
		})
	}

	return block{
		content: contentRows,
		options: options,
	}
}

func validOrPlaceholder(value any) string {
	if value == core.InvalidValue {
		return "-"
	}
	switch cast := value.(type) {
	case string:
		return cast
	case float64:
		return fmt.Sprintf("%.2f%%", value)
	case int:
		return fmt.Sprintf("%d", value)
	default:
		return fmt.Sprint(value)
	}
}
