// +build !codeanalysis

package graphs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdder(t *testing.T) {
	assert := assert.New(t)

	adder := NewAdder()
	sum := adder.Run(3, 2)
	assert.Equal(5, sum)
}
