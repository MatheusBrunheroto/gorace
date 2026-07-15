package verbose

import "math"

type rgb struct {
	red   int
	green int
	blue  int
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func hslToRGB(h, s, l float64) rgb {

	c := (1 - abs(2*l-1)) * s
	x := c * (1 - abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return rgb{
		red:   int((r + m) * 255),
		green: int((g + m) * 255),
		blue:  int((b + m) * 255),
	}

}

func hashToVividColor(hash uint64) rgb {
	hue := math.Mod(340+float64(hash%40), 360)
	lightness := 0.45 + float64((hash>>8)%25)/100
	return hslToRGB(hue, 0.85, lightness)
}
