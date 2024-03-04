package shared

import "image/color"

func GetWN8Color(r int) color.Color {
	if r > 0 && r < 301 {
		return color.RGBA{255, 0, 0, 255}
	}
	if r > 300 && r < 451 {
		return color.RGBA{251, 83, 83, 255}
	}
	if r > 450 && r < 651 {
		return color.RGBA{255, 160, 49, 255}
	}
	if r > 650 && r < 901 {
		return color.RGBA{255, 244, 65, 255}
	}
	if r > 900 && r < 1201 {
		return color.RGBA{149, 245, 62, 255}
	}
	if r > 1200 && r < 1601 {
		return color.RGBA{103, 190, 51, 255}
	}
	if r > 1600 && r < 2001 {
		return color.RGBA{106, 236, 255, 255}
	}
	if r > 2000 && r < 2451 {
		return color.RGBA{46, 174, 193, 255}
	}
	if r > 2450 && r < 2901 {
		return color.RGBA{208, 108, 255, 255}
	}
	if r > 2900 {
		return color.RGBA{142, 65, 177, 255}
	}
	return color.Transparent
}
