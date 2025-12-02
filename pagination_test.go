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

func TestBadToken(t *testing.T) {
	token := "asdfasdfasdf"
	_, err := pagination.Decode(token, []byte{})
	if !errors.Is(err, pagination.ErrInvalidToken) {
		t.Errorf("Decode(%v): unexpected error; want %v got %v", token, pagination.ErrInvalidToken, err)
	}
}

func ExampleDecode() {
	// Compute a hash of the request arguments.  AIP-158 says "the user is
	// expected to keep all other arguments to the RPC the same; if any
	// arguments are different, the API should send an INVALID_ARGUMENT error.
	paramHash := sha256.New()
	paramHash.Write([]byte("filter=an expression")) // from `req.params.filter`
	nonce := paramHash.Sum(nil)

	// In Request N-1

	pageOffset := 3 // current page offset from req params

	nextPageToken, err := pagination.Encode(pageOffset, nonce)
	if err != nil {
		fmt.Printf("InvalidArgument error: %v", err)
	}

	// resp.nextPageToken = nextPageToken
	// return resp, nil

	// In request N
	page, err := pagination.Decode(nextPageToken, nonce)
	if err != nil {
		// If the nonce differs, Decode will fail, warning us the request arguments differ.
		fmt.Printf("InvalidArgument error: %v", err)
	}

	fmt.Printf("Next page offset is %d", page)
	// Output: Next page offset is 4
}
