// Package pagination assists with RESTful collection listing APIs,
// implementing an [`AIP-158`](https://google.aip.dev/158)  compliant function
// to encode and decode `next_page_token`s into page sizes and offsets for
// continuing subsequent list method calls.
package pagination

// Decode takes a `next_page_token` UTF8 string and a nonce and returns a
// decoded next page number, or an error explaining why this token was not
// valid.  The nonce must be constant for the same request parameters, e.g. a
// hash of a filter expression string.  If the token is empty, a new Pagination
// starting at an offset of zero, with a page set to Size is returned.  Callers
// are responsible for converting the page number returned to an offset for
// their storage query engine.
func Decode(token string, nonce []byte, size int) (page int, err error) {
	return 0, nil
}

// Encode takes a page size and page number for the current response, and a
// nonce, and returns an encoded `next_page_token` UTF8 string to pass to a
// REST client as a way to continue a list query at the next page after this
// one.  If the encoding fails, an error is returned instead and the token is
// undefined.  The nonce must be constant for the same request parameters,
// e.g. a hash of the filter expression string.
func Encode(size, page int, nonce []byte) (next_page_token string, err error) {
	return "", nil
}
