package main

func max(in ...int) int {
	_max := in[0]
	for _, val := range in {
		if val > _max {
			_max = val
		}
	}
	return _max
}

func min(in ...int) int {
	_min := in[0]
	for _, val := range in {
		if val < _min {
			_min = val
		}
	}
	return _min
}

func sum(in ...int) int {
	_sum := 0
	for _, val := range in {
		_sum += val
	}
	return _sum
}

func isValidHole(hole int) bool {
	if hole > 7 || hole < 1 {
		return false
	}
	return true
}

func maxf(in ...float64) float64 {
	_max := in[0]
	for _, val := range in {
		if val > _max {
			_max = val
		}
	}
	return _max
}

func minf(in ...float64) float64 {
	_min := in[0]
	for _, val := range in {
		if val < _min {
			_min = val
		}
	}
	return _min
}

func sumf(in ...float64) float64 {
	var _sum float64
	_sum = 0.0
	for _, val := range in {
		_sum += val
	}
	return _sum
}
