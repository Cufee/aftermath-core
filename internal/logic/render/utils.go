package render

import (
	"fmt"
	"image"

	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/fogleman/gg"
)

func renderImages(images []image.Image, style Style) image.Image {
	if len(images) == 0 {
		return nil
	}

	totalWidth := 0
	totalHeight := 0
	maxWidth := 0
	maxHeight := 0
	for _, img := range images {
		totalWidth += img.Bounds().Dx()
		if img.Bounds().Dx() > maxWidth {
			maxWidth = img.Bounds().Dx()
		}
		totalHeight += img.Bounds().Dy()
		if img.Bounds().Dy() > maxHeight {
			maxHeight = img.Bounds().Dy()
		}
	}

	imageWidth := maxWidth
	imageHeight := maxHeight
	if style.Direction == DirectionHorizontal {
		imageWidth = totalWidth + (len(images)-1)*int(style.Gap)
	} else {
		imageHeight = totalHeight + (len(images)-1)*int(style.Gap)
	}

	ctx := gg.NewContext(imageWidth, imageHeight)
	var lastX, lastY float64
	for i, img := range images {
		posX, posY := lastX, lastY

		switch style.Direction {
		case DirectionHorizontal:
			if i > 0 {
				posX += style.Gap
			}
			lastX = posX + float64(img.Bounds().Dx())

			switch style.AlignItems {
			case AlignItemsCenter:
				posY = float64(imageHeight-img.Bounds().Dy()) / 2
			case AlignItemsEnd:
				posY = float64(imageHeight - img.Bounds().Dy())
			}

		case DirectionVertical:
			if i > 0 {
				posY += style.Gap
			}
			lastY = posY + float64(img.Bounds().Dy())

			switch style.AlignItems {
			case AlignItemsCenter:
				posX = float64(imageWidth-img.Bounds().Dx()) / 2
			case AlignItemsEnd:
				posX = float64(imageWidth - img.Bounds().Dx())
			}
		}

		ctx.DrawImage(img, int(posX), int(posY))
	}

	return ctx.Image()
}

func newDataRow(value any, config RenderConfig) blockRow {
	return blockRow{
		value:  validOrPlaceholder(value),
		config: config,
	}
}

func newLabelRow(label blockLabelTag, locale localization.SupportedLanguage, config RenderConfig) blockRow {
	if string(locale) == "" {
		locale = localization.LanguageEN
	}
	return blockRow{
		value:  string(label),
		locale: &locale,
		config: config,
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
