package worktile

import (
	"fmt"
	"testing"
)

func TestEmtpySlice(t *testing.T) {
	s := []string{"a", "b", "c"}
	fmt.Println(s, len(s), cap(s))
	s = nil
	fmt.Println(s, len(s), cap(s))
	s = append(s, "d")
	fmt.Println(s, len(s), cap(s))
}