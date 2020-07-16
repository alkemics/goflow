package functions

type IntReducer func(list ...int) int

var IntSum IntReducer = func(list ...int) int {
	sum := 0
	for _, e := range list {
		sum += e
	}
	return sum
}

var IntMultiplication IntReducer = func(list ...int) int {
	res := 1
	for _, e := range list {
		res *= e
	}
	return res
}

type UIntReducer func(list ...uint) uint

var UIntSum UIntReducer = func(list ...uint) uint {
	var sum uint = 0
	for _, e := range list {
		sum += e
	}
	return sum
}
