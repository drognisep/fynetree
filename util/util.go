package util

import "fyne.io/fyne"

func IntMax(ints ...int) int {
	var max int
	for i, num := range ints {
		if i == 0 {
			max = num
			continue
		}
		if num > max {
			max = num
		}
	}
	return max
}

func InlineMinSize(sizes ...fyne.Size) fyne.Size {
	var runningWidth int
	var maxHeight int
	for _, size := range sizes {
		runningWidth += size.Width
		maxHeight = IntMax(size.Height, maxHeight)
	}
	return fyne.Size{
		Width:  runningWidth,
		Height: maxHeight,
	}
}

func ColumnMinSize(sizes ...fyne.Size) fyne.Size {
	var maxWidth int
	var runningHeight int
	for _, size := range sizes {
		maxWidth = IntMax(size.Width, maxWidth)
		runningHeight += size.Height
	}
	return fyne.Size{
		Width:  maxWidth,
		Height: runningHeight,
	}
}
