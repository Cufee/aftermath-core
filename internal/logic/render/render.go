package render

import (
	"errors"
	"image"
	"math"

	"github.com/fogleman/gg"
)

type RenderOptions struct {
	Style Style
	Debug bool
}

func renderImages(images []image.Image, options *RenderOptions) (image.Image, error) {
	if len(images) == 0 {
		return nil, errors.New("no images to render")
	}

	imageSize := getDetailedSize(images, options.Style)

	var lastX, lastY float64 = options.Style.PaddingX, options.Style.PaddingY
	var justifyOffsetX, justifyOffsetY float64
	var elementWidth, elementHeight float64

	/*
		TODO: Some math here under some configurations is certainly broken and needs to be fixed or removed
	*/

	// Set correct gaps and offsets based on justify content
	switch options.Style.JustifyContent {
	case JustifyContentCenter:
		lastX += float64(imageSize.extraSpacingX / 2)
		lastY += float64(imageSize.extraSpacingY / 2)
	case JustifyContentEnd:
		lastX += float64(imageSize.extraSpacingX)
		lastY += float64(imageSize.extraSpacingY)
	case JustifyContentSpaceBetween:
		justifyOffsetX = float64(imageSize.extraSpacingX / float64(len(images)-1))
		justifyOffsetY = float64(imageSize.extraSpacingY / float64(len(images)-1))
	case JustifyContentSpaceEvenly:
		elementWidth = imageSize.maxElementWidth
		elementHeight = imageSize.maxElementHeight
		// All spacing is equal, so we can just use the same formula as space around
		fallthrough
	case JustifyContentSpaceAround:
		spacingX := float64(imageSize.extraSpacingX / float64(len(images)+1))
		spacingY := float64(imageSize.extraSpacingY / float64(len(images)+1))
		lastX += spacingX
		lastY += spacingY
		justifyOffsetX = spacingX
		justifyOffsetY = spacingY
	default: // JustifyContentStart
		// 0,0
	}

	ctx := gg.NewContext(int(math.Ceil(imageSize.width)), int(math.Ceil(imageSize.height)))

	if options.Style.BorderRadius > 0 {
		ctx.DrawRoundedRectangle(0, 0, float64(ctx.Width()), float64(ctx.Height()), options.Style.BorderRadius)
		ctx.Clip()
	}
	if options.Style.BackgroundColor != nil {
		ctx.DrawRectangle(0, 0, imageSize.width, imageSize.height)
		ctx.SetColor(options.Style.BackgroundColor)
		ctx.Fill()
	}

	for i, img := range images {
		posX, posY := lastX, lastY

		imgWidth := float64(img.Bounds().Dx())
		imgHeight := float64(img.Bounds().Dy())

		targetWidth := max(imgWidth, elementWidth)
		targetHeight := max(imgHeight, elementHeight)

		switch options.Style.Direction {
		case DirectionVertical:
			if i > 0 {
				posY += max(options.Style.Gap, justifyOffsetY)
			}
			lastY = posY + targetHeight

			switch options.Style.AlignItems {
			case AlignItemsCenter:
				posX = (imageSize.width - imgWidth) / 2
			case AlignItemsEnd:
				posX = imageSize.height - imgWidth
			default: // AlignItemsStart
				posX = 0
			}
		default: // DirectionHorizontal
			if i > 0 {
				posX += max(options.Style.Gap, justifyOffsetX)
			}
			lastX = posX + targetWidth

			switch options.Style.AlignItems {
			case AlignItemsCenter:
				posY = (imageSize.height - imgHeight) / 2
			case AlignItemsEnd:
				posY = imageSize.height - imgHeight
			default: // AlignItemsStart
				posY = 0
			}

		}

		if options.Debug {
			ctx.SetColor(debugColorPink)
			ctx.DrawRectangle(posX, posY, targetWidth, targetHeight)
			ctx.Stroke()
		}

		ctx.DrawImage(img, int(math.Ceil(posX+(targetWidth-imgWidth)/2)), int(math.Ceil(posY+(targetHeight-imgHeight)/2)))
	}

	return ctx.Image(), nil
}

type imageSize struct {
	width  float64
	height float64

	// The amount of extra spacing added to the image, used for alignment
	extraSpacingX float64
	extraSpacingY float64

	totalGap float64

	maxElementWidth  float64
	maxElementHeight float64
}

func getDetailedSize(images []image.Image, style Style) imageSize {
	imageWidth, imageHeight := style.Width, style.Height

	var totalGap float64
	if len(images) > 1 {
		totalGap = float64(len(images)-1) * style.Gap
	}

	var totalWidth float64 = style.PaddingX * 2
	var totalHeight float64 = style.PaddingY * 2

	maxWidth, maxHeight := 0.0, 0.0

	for _, img := range images {
		imgX := float64(img.Bounds().Dx())
		if imgX > maxWidth {
			maxWidth = imgX
		}

		imgY := float64(img.Bounds().Dy())
		if imgY > maxHeight {
			maxHeight = imgY
		}

		if style.Direction == DirectionHorizontal {
			totalWidth += float64(img.Bounds().Dx())
		} else {
			totalHeight += float64(img.Bounds().Dy())
		}
	}

	if style.Width == 0 {
		imageWidth = totalWidth
	}
	if style.Height == 0 {
		imageHeight = totalHeight
	}

	extraSpacingX := imageWidth - totalWidth
	extraSpacingY := imageHeight - totalHeight

	switch style.Direction {
	case DirectionVertical:
		if style.JustifyContent == JustifyContentSpaceEvenly && style.Height == 0 {
			imageHeight = totalGap + (style.PaddingY * 2) + (maxHeight * float64(len(images)))
		} else {
			imageHeight += totalGap
		}
		if style.Width == 0 {
			imageWidth = maxWidth + (style.PaddingX * 2)
		}
	default: // DirectionHorizontal
		if style.JustifyContent == JustifyContentSpaceEvenly && style.Width == 0 {
			imageWidth = totalGap + (style.PaddingX * 2) + (maxWidth * float64(len(images)))
		} else {
			imageWidth += totalGap
		}
		if style.Height == 0 {
			imageHeight = maxHeight + (style.PaddingY)*2
		}
	}

	return imageSize{
		totalGap:         totalGap,
		width:            imageWidth,
		height:           imageHeight,
		extraSpacingX:    extraSpacingX,
		extraSpacingY:    extraSpacingY,
		maxElementWidth:  maxWidth,
		maxElementHeight: maxHeight,
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
