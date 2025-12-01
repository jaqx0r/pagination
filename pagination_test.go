package pagination_test

import (
	"errors"
	"testing"

	"github.com/jaqx0r/pagination"
)

func TestNotNegative(t *testing.T) {
	_, err := pagination.Encode(-1, 0, []byte{})
	if !errors.Is(err, pagination.ErrNegativePageSize) {
		t.Errorf("Encode(-1): expected ErrNegativePageSize, got %v", err)
	}
}
