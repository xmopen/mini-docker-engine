package random

import (
	"fmt"
	"testing"
)

func TestRandNumberToString(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(RandNumberToString(10))
	}
}
