package pagination_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jaqx0r/pagination"
)

func TestRoundTrip(t *testing.T) {
	token, err := pagination.Encode(0, []byte{})
	if err != nil {
		t.Errorf("Encode(0): %v", err)
	}

	page, err := pagination.Decode(token, []byte{})
	if err != nil {
		t.Errorf("Decode(%v): %v", token, err)
	}
	if page != 1 {
		t.Errorf("unexpected page number, expctint 1, got %d", page)
	}
}

func TestChangedParameters(t *testing.T) {
	token, err := pagination.Encode(0, []byte("a"))
	if err != nil {
		t.Fatalf("Encode(0): %v", err)
	}

	_, err = pagination.Decode(token, []byte("b"))
	if !cmp.Equal(err, pagination.ErrChangedParameters, cmpopts.EquateErrors()) {
		t.Errorf("Decode(%v): unexpected error; want %v got %v", token, pagination.ErrChangedParameters, err)
	}

}
