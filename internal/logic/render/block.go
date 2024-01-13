package render

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/fogleman/gg"
)

type blockLabelTag string

const (
	blockLabelTagBattles  blockLabelTag = "battle"
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
	rows    []blockRow
	options RenderOptions
}

func (block *block) Render() (image.Image, error) {
	var images []image.Image
	for _, row := range block.rows {
		img, err := row.Render()
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return renderImages(images, block.options)
}

type blockRow struct {
	value  string
	locale *localization.SupportedLanguage

	config RenderConfig
}

func (row *blockRow) Render() (image.Image, error) {
	if row.config.Font == nil {
		return nil, errors.New("font not set")
	}

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
	ctx := gg.NewContext(int(math.Ceil(valueW)+1), int(math.Ceil(valueH+(descenderOffset*2))))

	// Render text
	ctx.SetFontFace(row.config.Font)
	ctx.SetColor(row.config.FontColor)

	ctx.DrawString(stringValue, 0, valueH)

	return ctx.Image(), nil
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
		rows = append(rows, cfg.NewLabel(label, cfg.Locale))
	}

	return block{
		rows:    rows,
		options: cfg.RowOptions,
	}
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
