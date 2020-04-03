package array

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInArrayInt(t *testing.T) {
	arr := []int{1, 2, 3}
	assert.True(t, InArrayInt(arr, 1))
	assert.False(t, InArrayInt(arr, 6))
}

func TestRemoveDuplicate(t *testing.T) {
	arr := []string{"A", "A", "B", "C", "C"}
	expected := []string{"A", "B", "C"}
	actual := RemoveDuplicate(arr)
	sort.Strings(actual)
	assert.Equal(t, expected, actual)
}

func TestInArrayString(t *testing.T) {
	arr := []string{"1", "2", "3"}
	assert.True(t, InArrayString(arr, "1"))
	assert.False(t, InArrayString(arr, "6"))
}
