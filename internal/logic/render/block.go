package render

import (
	"image"
	"image/color"
	"math"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

type blockLabelTag string

const (
	blockLabelTagBattles  blockLabelTag = "battles"
	blockLabelTagAvgDmg   blockLabelTag = "avgDmg"
	blockLabelTagWinrate  blockLabelTag = "winrate"
	blockLabelTagAccuracy blockLabelTag = "accuracy"
	blockLabelTagWN8      blockLabelTag = "wn8"

	blockLabelTagNone blockLabelTag = "none"
)

var (
	debugColorPink = color.RGBA{255, 192, 203, 255}
)

type block struct {
	rows  []blockRow
	style Style
}

func (block *block) Render(style Style) image.Image {
	var images []image.Image
	for _, row := range block.rows {
		images = append(images, row.Render())
	}
	return renderImages(images, style)
}

type blockSet struct {
	blocks []block
	style  Style
}

func (set *blockSet) Render() image.Image {
	var images []image.Image
	for _, block := range set.blocks {
		images = append(images, block.Render(block.style))
	}
	return renderImages(images, set.style)
}

type blockRow struct {
	value  string
	locale *localization.SupportedLanguage

	config RenderConfig
}

func (row *blockRow) Render() image.Image {
	stringValue := row.value
	if row.locale != nil {
		// TODO: Implement localization
		_ = row.locale
	}

	measureCtx := gg.NewContext(1, 1)
	measureCtx.SetFontFace(row.config.Font)
	valueW, valueH := measureCtx.MeasureString(stringValue)

	// Account for font descender height
	descenderOffset := (float64(row.config.Font.Metrics().Descent>>6) - 1)
	ctx := gg.NewContext(int(math.Ceil(valueW)), int(math.Ceil(valueH+(descenderOffset))))

	if row.config.Debug {
		ctx.SetColor(debugColorPink)
		ctx.Clear()
		ctx.SetColor(color.Black)
	}

	// Render text
	ctx.SetFontFace(row.config.Font)
	ctx.SetColor(row.config.FontColor)

	ctx.DrawString(stringValue, -1, valueH-descenderOffset)

	return ctx.Image()
}

type RenderConfig struct {
	Font      font.Face
	FontColor color.RGBA

	Debug bool
}

type BlockRenderConfig struct {
	Session RenderConfig `json:"session"`
	Career  RenderConfig `json:"career"`
	Label   RenderConfig `json:"label"`

	RowStyle Style `json:"rowStyle"`
	SetStyle Style `json:"setStyle"`
}

func (cfg *BlockRenderConfig) NewLabel(label blockLabelTag, locale localization.SupportedLanguage) blockRow {
	return newLabelRow(label, locale, cfg.Label)
}

func (cfg *BlockRenderConfig) NewSessionRow(value any) blockRow {
	return newDataRow(value, cfg.Session)
}

func (cfg *BlockRenderConfig) NewCareerRow(value any) blockRow {
	return newDataRow(value, cfg.Career)
}

func (cfg *BlockRenderConfig) CompleteBlock(label blockLabelTag, session, career any) block {
	var rows []blockRow
	if session != nil {
		rows = append(rows, cfg.NewSessionRow(session))
	}
	if career != nil {
		rows = append(rows, cfg.NewCareerRow(career))
	}
	if label != blockLabelTagNone {
		rows = append(rows, cfg.NewLabel(label, localization.LanguageEN))
	}

	return block{
		rows:  rows,
		style: cfg.RowStyle,
	}
}
