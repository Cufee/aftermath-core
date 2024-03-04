package shared

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type LogoSizingOptions struct {
	Lines int

	LineStep  int
	LineWidth float64
	Jump      float64
	Gap       float64
}

func (opts LogoSizingOptions) Height() int {
	return ((opts.Lines/2+1)*opts.LineStep + (opts.Lines/2)*int(opts.Jump))
}
func (opts LogoSizingOptions) Width() int {
	return opts.Lines * (int(opts.LineWidth + opts.Gap))
}

func DefaultLogoOptions() LogoSizingOptions {
	return LogoSizingOptions{
		Gap:       4,
		Jump:      6,
		Lines:     9,
		LineStep:  12,
		LineWidth: 6,
	}
}

func AftermathLogo(fillColor color.Color, opts LogoSizingOptions) image.Image {
	ctx := gg.NewContext(opts.Width(), opts.Height())
	for line := range opts.Lines {
		height := opts.LineStep + opts.LineStep*line

		offset := float64(opts.Height() - height - (line * int(opts.Jump)))
		if line > opts.Lines/2 {
			height = opts.LineStep + opts.LineStep*(opts.Lines-line-1)
			offset = float64(opts.Height() - height - ((opts.Lines - line - 1) * int(opts.Jump)))
		}

		ctx.DrawRoundedRectangle((opts.Gap/2)+float64(line*(int(opts.LineWidth+opts.Gap))), offset, opts.LineWidth, float64(height), 3)
		ctx.SetColor(fillColor)
		ctx.Fill()
		ctx.ClearPath()
	}

	return ctx.Image()
}
