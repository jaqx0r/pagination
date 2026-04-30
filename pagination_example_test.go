package pagination_test

import (
	"crypto/sha256"
	"fmt"

	"github.com/jaqx0r/pagination"
)

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

	cursor, err := pagination.Decode[int](req.pageToken, nonce)
	if err != nil {
		fmt.Printf("Invalid argument: %v\n", err)
		return
	}

	// Perform query
	fmt.Printf("Query with cursor offset: %d and limit: %d\n", cursor, req.pageSize)

	nextPageToken, err := pagination.Encode[int](cursor+req.pageSize, nonce)
	if err != nil {
		fmt.Printf("Internal server error: %v\n", err)
	}

	// response.nextPageToken = nextPageToken
	req.pageToken = nextPageToken
	// return resp, nil

	// In request 2
	// recompute nonce
	cursor2, err := pagination.Decode[int](req.pageToken, nonce)
	if err != nil {
		fmt.Printf("Invalid argument: %v\n", err)
		return
	}

	fmt.Printf("Next page offset is %d\n", cursor2)
	// Output:
	// Query with cursor offset: 0 and limit: 10
	// Next page offset is 10
}
