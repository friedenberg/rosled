package main

const (
	Rows       = 36
	Columns    = 96
	PixelCount = 3648 // Rows * Columns
)

type Display [PixelCount]bool
