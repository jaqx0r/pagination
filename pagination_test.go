package pagination_test

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jaqx0r/pagination"
)

func TestRoundTrip(t *testing.T) {
	token, err := pagination.Encode(0, 1, []byte{})
	if err != nil {
		t.Errorf("Encode(0, 0): %v", err)
	}

	offset, err := pagination.Decode(token, []byte{})
	if err != nil {
		t.Errorf("Decode(%v): %v", token, err)
	}
	if offset != 1 {
		t.Errorf("unexpected page number, expecting 1, got %d", offset)
	}
}

func TestChangedPageSize(t *testing.T) {
	token, err := pagination.Encode(0, 1, []byte{})
	if err != nil {
		t.Errorf("Encode(0, 0): %v", err)
	}

	offset, err := pagination.Decode(token, []byte{})
	if err != nil {
		t.Errorf("Decode(%v): %v", token, err)
	}

	token2, err := pagination.Encode(offset, 10, []byte{})
	if err != nil {
		t.Errorf("Encode(0, 0): %v", err)
	}

	offset2, err := pagination.Decode(token2, []byte{})
	if err != nil {
		t.Errorf("Decode(%v): %v", token, err)
	}
	if offset2 != 11 {
		t.Errorf("unexpected page number, expecting 11, got %d", offset2)
	}
}

func TestChangedParameters(t *testing.T) {
	token, err := pagination.Encode(0, 0, []byte("a"))
	if err != nil {
		t.Fatalf("Encode(0): %v", err)
	}

	_, err = pagination.Decode(token, []byte("b"))
	if !cmp.Equal(err, pagination.ErrChangedParameters, cmpopts.EquateErrors()) {
		t.Errorf("Decode(%v): unexpected error; want %v got %v", token, pagination.ErrChangedParameters, err)
	}
}

func TestBadToken(t *testing.T) {
	token := "asdfasdfasdf"
	_, err := pagination.Decode(token, []byte{})
	if !errors.Is(err, pagination.ErrInvalidToken) {
		t.Errorf("Decode(%v): unexpected error; want %v got %v", token, pagination.ErrInvalidToken, err)
	}
}

func TestEmptyToken(t *testing.T) {
	page, err := pagination.Decode("", []byte{})
	if err != nil {
		t.Errorf("Decode('') unexpected error: %v", err)
	}
	if page != 0 {
		t.Errorf("page: want %v got %v", 0, page)
	}
}

func ExampleDecode() {
	// simulate request input
	req := struct {
		pageToken string
		filter    string
		pageSize  int
	}{
		pageToken: "",
		filter:    "a query expression",
		pageSize:  10,
	}

	// Request 1
	paramSignature := sha256.New()
	paramSignature.Write([]byte(req.filter))
	nonce := paramSignature.Sum(nil)

	offset, err := pagination.Decode(req.pageToken, nonce)
	if err != nil {
		fmt.Printf("Invalid argument: %v\n", err)
		return
	}

	// Perform query
	fmt.Printf("Query with cursor offset: %d and limit: %d\n", offset, req.pageSize)

	nextPageToken, err := pagination.Encode(offset, req.pageSize, nonce)
	if err != nil {
		fmt.Printf("Internal server error: %v\n", err)
	}

	// response.nextPageToken = nextPageToken
	req.pageToken = nextPageToken
	// return resp, nil

	// In request 2
	// recompute nonce
	offset2, err := pagination.Decode(req.pageToken, nonce)
	if err != nil {
		fmt.Printf("Invalid argument: %v\n", err)
		return
	}

	fmt.Printf("Next page offset is %d\n", offset2)
	// Output:
	// Query with cursor offset: 0 and limit: 10
	// Next page offset is 10
}
