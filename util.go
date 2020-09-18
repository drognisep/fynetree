package fynetree

func intMax(ints ...int) int {
	var max int
	for _, i := range ints {
		if i > max {
			max = i
		}
	}
	return max
}
