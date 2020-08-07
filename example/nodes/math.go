package nodes

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alkemics/goflow/example/constants/functions"
)

func Add(a, b int) (sum int) { return a + b }

func Adder(a, b int) (sum int) { return a + b }

func Multiplier(a, b int) (product int) { return a * b }

func RandomIntProducer() (n int) {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100) //nolint:gosec // not really critical
}

// IntAggregator aggregates `list` using `reducer.
// optional: list
func IntAggregator(list []int, reducer functions.IntReducer) (result int) {
	return reducer(list...)
}

func UAdd(a, b uint) (sum uint) { return a + b }

func UIntAggregator(list []uint, reducer functions.UIntReducer) (result uint) {
	return reducer(list...)
}

type IntReducer struct {
	debug bool
}

func NewIntReducer(debug bool) *IntReducer {
	return &IntReducer{debug: debug}
}

func (r IntReducer) reduce(list []int, reducer functions.IntReducer) int {
	result := reducer(list...)
	if r.debug {
		fmt.Println(result)
	}
	return result
}

// Add adds a list of ints.
func (r IntReducer) Add(list []int) (sum int) {
	return r.reduce(list, functions.IntSum)
}

func (r IntReducer) Multiply(list []int) (product int) {
	return r.reduce(list, functions.IntMultiplication)
}
