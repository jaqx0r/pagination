package pagination_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jaqx0r/pagination"
)

type roundTripTest[T any] struct {
	name   string
	cursor T
}

func (rtt roundTripTest[T]) Name() string {
	return rtt.name
}

func (rtt roundTripTest[T]) Test(t *testing.T) {
	token, err := pagination.Encode(rtt.cursor, []byte{})
	if err != nil {
		t.Fatalf("Encode(%v): %v", rtt.cursor, err)
	}
	cursor, err := pagination.Decode[T](token, []byte{})
	if err != nil {
		t.Errorf("Decode(%v): %v", token, err)
	}
	if diff := cmp.Diff(rtt.cursor, cursor); diff != "" {
		t.Fatalf("Decode(%v): unexpected cursor returned (-want +got):\n%s", token, diff)
	}
}

type roundTripTestable interface {
	Test(t *testing.T)
	Name() string
}

type LimitOffset struct {
	Limit  int
	Offset int
}

type IDCursor struct {
	ID   string
	Name string
}

func TestRoundTrip(t *testing.T) {
	for _, tc := range []roundTripTestable{
		roundTripTest[string]{"string", "test"},
		roundTripTest[LimitOffset]{"LimitOffset", LimitOffset{0, 1}},
		roundTripTest[IDCursor]{"IDCursor", IDCursor{"1", "test"}},
	} {
		t.Run(tc.Name(), tc.Test)
	}
}

func TestChangedParameters(t *testing.T) {
	token, err := pagination.Encode(0, []byte("a"))
	if err != nil {
		t.Fatalf("Encode(0): %v", err)
	}

	_, err = pagination.Decode[int](token, []byte("b"))
	if !cmp.Equal(err, pagination.ErrChangedParameters, cmpopts.EquateErrors()) {
		t.Errorf("Decode(%v): unexpected error; want %v got %v", token, pagination.ErrChangedParameters, err)
	}
}

func TestBadToken(t *testing.T) {
	token := "asdfasdfasdf"
	_, err := pagination.Decode[int](token, []byte{})
	if !errors.Is(err, pagination.ErrInvalidToken) {
		t.Errorf("Decode(%v): unexpected error; want %v got %v", token, pagination.ErrInvalidToken, err)
	}
}

func TestEmptyToken(t *testing.T) {
	page, err := pagination.Decode[int]("", []byte{})
	if err != nil {
		t.Errorf("Decode('') unexpected error: %v", err)
	}
	if page != 0 {
		t.Errorf("page: want %v got %v", 0, page)
	}
}
