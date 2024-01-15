package render

import (
	"errors"
	"image"
	"math"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

type blockContentType int

const (
	BlockContentTypeText blockContentType = iota
	BlockContentTypeImage
	// BlockContentTypeIcon
	BlockContentTypeBlocks
)

type BlockContent interface {
	Render(Style) (image.Image, error)
	Type() blockContentType
}

type Block struct {
	ContentType blockContentType
	content     BlockContent
	Style       Style
}

func (block *Block) Render() (image.Image, error) {
	return block.content.Render(block.Style)
}

func NewBlock(content BlockContent, style Style) Block {
	return Block{
		ContentType: content.Type(),
		content:     content,
		Style:       style,
	}
}

type contentText struct {
	value string
}

func NewTextContent(style Style, value string) Block {
	return NewBlock(contentText{
		value: value,
	}, style)
}

func (content contentText) Render(style Style) (image.Image, error) {
	if style.Font == nil {
		return nil, errors.New("font not set")
	}

	measureCtx := gg.NewContext(1, 1)
	measureCtx.SetFontFace(*style.Font)
	valueW, valueH := measureCtx.MeasureString(content.value)

	// Account for font descender height
	descenderOffset := (float64((*style.Font).Metrics().Descent>>6) - 1)
	ctx := gg.NewContext(int(style.PaddingX*2+math.Ceil(valueW)+1), int(style.PaddingY*2+math.Ceil(valueH+(descenderOffset*2))))

	// Render text
	ctx.SetFontFace(*style.Font)
	ctx.SetColor(style.FontColor)

	ctx.DrawString(content.value, style.PaddingX, valueH+style.PaddingY)

	if style.Debug {
		ctx.SetColor(getDebugColor())
		ctx.DrawRectangle(0, 0, float64(ctx.Width()), float64(ctx.Height()))
		ctx.Stroke()
	}

	return ctx.Image(), nil
}

func (content contentText) Type() blockContentType {
	return BlockContentTypeText
}

type contentBlocks struct {
	blocks []Block
}

func NewBlocksContent(style Style, blocks ...Block) Block {
	return NewBlock(contentBlocks{
		blocks: blocks,
	}, style)
}

func (content contentBlocks) Render(style Style) (image.Image, error) {
	var images []image.Image
	for _, block := range content.blocks {
		img, err := block.Render()
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return renderImages(images, style)
}

func (content contentBlocks) Type() blockContentType {
	return BlockContentTypeBlocks
}

type contentImage struct {
	image image.Image
}

func NewImageContent(style Style, image image.Image) Block {
	return NewBlock(contentImage{
		image: image,
	}, style)
}

func (content contentImage) Render(style Style) (image.Image, error) {
	if style.Width == 0 || style.Height == 0 {
		return nil, errors.New("width or height not set")
	}

	// Type cast to image.Image for gg
	var image image.Image = imaging.Fit(content.image, int(style.Width), int(style.Height), imaging.Linear)
	if style.BackgroundColor != nil {
		mask := gg.NewContextForImage(image)
		ctx := gg.NewContext(image.Bounds().Dx(), image.Bounds().Dy())
		ctx.SetMask(mask.AsMask())
		ctx.SetColor(style.BackgroundColor)
		ctx.DrawRectangle(0, 0, float64(ctx.Width()), float64(ctx.Height()))
		ctx.Fill()
		image = ctx.Image()
	}

	return image, nil
}

func (content contentImage) Type() blockContentType {
	return BlockContentTypeBlocks
}
