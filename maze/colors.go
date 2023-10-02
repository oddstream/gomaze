package maze

import (
	"image/color"
)

// BasicColors is a string-indexed map of the basic colors from HTML 4.01
var BasicColors = map[string]color.RGBA{
	"White":  {R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	"Silver": {R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff},
	"Gray":   {R: 0x80, G: 0x80, B: 0x80, A: 0xff},
	"Black":  {R: 0x00, G: 0x00, B: 0x00, A: 0xff},
	"Red":    {R: 0xff, G: 0, B: 0, A: 0xff},
	"Maroon": {R: 0x80, G: 0, B: 0, A: 0xff},
	"Yellow": {R: 0xff, G: 0xff, B: 0, A: 0xff},
	"Olive":  {R: 0x80, G: 0x80, B: 0, A: 0xff},
	"Lime":   {R: 0, G: 0xff, B: 0, A: 0xff},
	"Green":  {R: 0, G: 0x80, B: 0, A: 0xff},
	"Aqua":   {R: 0, G: 0xff, B: 0xff, A: 0xff},
	"Teal":   {R: 0, G: 0x80, B: 0x80, A: 0xff},
	"Blue":   {R: 0, G: 0, B: 0xff, A: 0xff},
	"Navy":   {R: 0, G: 0, B: 0x80, A: 0xff},
	"Fushia": {R: 0xff, G: 0, B: 0xff, A: 0xff},
	"Purple": {R: 0x80, G: 0, B: 0x80, A: 0xff},
}

// ExtendedColors see en.wikipedia.org/wiki/Web_colors
var ExtendedColors = map[string]color.RGBA{
	// Pink colors
	"MediumVioletRed": {R: 0xc7, G: 0x15, B: 0x85, A: 0xff},
	"DeepPink":        {R: 0xff, G: 0x14, B: 0x93, A: 0xff},
	"PaleVioletRed":   {R: 0xdb, G: 0x70, B: 0x93, A: 0xff},
	"HotPink":         {R: 0xff, G: 0x69, B: 0xb4, A: 0xff},
	"LightPink":       {R: 0xff, G: 0xb6, B: 0xc1, A: 0xff},
	"Pink":            {R: 0xff, G: 0xc0, B: 0xcb, A: 0xff},

	"Navy":           {R: 0, G: 0, B: 0x80, A: 0xff},
	"DarkBlue":       {R: 0, G: 0, B: 0x8b, A: 0xff},
	"MediumBlue":     {R: 0, G: 0, B: 0xcd, A: 0xff},
	"Blue":           {R: 0, G: 0, B: 0xff, A: 0xff}, // glows horribly, do not use
	"MidnightBlue":   {R: 0x19, G: 0x19, B: 0x70, A: 0xff},
	"RoyalBlue":      {R: 0x41, G: 0x69, B: 0xe1, A: 0xff}, // glows horribly, do not use
	"SteelBlue":      {R: 0x46, G: 0x82, B: 0xb4, A: 0xff},
	"DodgerBlue":     {R: 0x1e, G: 0x90, B: 0xff, A: 0xff},
	"DeepSkyBlue":    {R: 0x0, G: 0xbf, B: 0xff, A: 0xff},
	"CornflowerBlue": {R: 0x64, G: 0x95, B: 0xed, A: 0xff},
	"SkyBlue":        {R: 0x87, G: 0xce, B: 0xeb, A: 0xff},
	"LightSkyBlue":   {R: 0x87, G: 0xce, B: 0xfa, A: 0xff},
	"LightSteelBlue": {R: 0xb0, G: 0xc4, B: 0xde, A: 0xff},
	"LightBlue":      {R: 0xad, G: 0xd8, B: 0xe6, A: 0xff},
	"PowderBlue":     {R: 0xb0, G: 0xe0, B: 0xe6, A: 0xff},

	"DarkKhaki":     {R: 0xbd, G: 0xb7, B: 0x6b, A: 0xff},
	"Gold":          {R: 0xff, G: 0xd7, B: 0, A: 0xff},
	"Khaki":         {R: 0xf0, G: 0xe6, B: 0x8c, A: 0xff},
	"PeachPuff":     {R: 0xff, G: 0xda, B: 0xb9, A: 0xff},
	"Yellow":        {R: 0xff, G: 0xff, B: 0, A: 0xff},
	"PaleGoldenRod": {R: 0xee, G: 0xe8, B: 0xaa, A: 0xff},
	"Moccasin":      {R: 0xff, G: 0xe4, B: 0xb5, A: 0xff},
	"PapayaWhip":    {R: 0xff, G: 0xef, B: 0xd5, A: 0xff},
} // golang gotcha no newline after last literal, must be comma or closing brace

// Palette a slice of color names as strings
type Palette = []string

// Palettes a slice of Palette
var Palettes = []Palette{
	{"SteelBlue", "CornflowerBlue", "SkyBlue", "LightSteelBlue", "LightBlue", "PowderBlue"},
	// {"MediumVioletRed", "DeepPink", "PaleVioletRed", "HotPink", "LightPink", "Pink"},
	{"Gold", "Khaki", "PeachPuff", "PaleGoldenRod", "Moccasin", "PapayaWhip", "DarkKhaki", "Yellow"},
}

var (
	colorBackground = color.RGBA{R: 0x50, G: 0x50, B: 0x50, A: 0xff}
)

// CalcBackgroundColor returns the average of all the colors in a palette, divided by two
func CalcBackgroundColor(colNames []string) color.RGBA {
	var r, g, b int
	for _, cnam := range colNames {
		col := ExtendedColors[cnam]
		r += int(col.R)
		g += int(col.G)
		b += int(col.B)
	}
	r = r / len(colNames) / 2
	g = g / len(colNames) / 2
	b = b / len(colNames) / 2
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}
}
