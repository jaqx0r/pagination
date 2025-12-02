// Package pagination assists with RESTful collection listing APIs,
// implementing an [`AIP-158`](https://google.aip.dev/158)  compliant function
// to encode and decode `next_page_token`s into page sizes and offsets for
// continuing subsequent list method calls.
package pagination

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
)

var ErrChangedParameters = errors.New("parameters changed between method calls invalidating this token")
var ErrInvalidToken = errors.New("invalid token")

type nextPageToken struct {
	Page  int
	Nonce []byte
}

// Decode takes a `next_page_token` UTF8 string and a nonce and returns a
// decoded next page number, or an error explaining why this token was not
// valid.  The nonce must be constant for the same query parameters, e.g. a
// hash of a filter expression string.  If the token is empty, the page is
// zero.  Callers are responsible for converting the page number returned to an
// offset for their storage query engine.
func Decode(token string, nonce []byte) (page int, err error) {
	if token == "" {
		return 0, nil
	}
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	bs := bytes.NewBuffer(b)

	var nextPageToken nextPageToken
	gDec := gob.NewDecoder(bs)
	err = gDec.Decode(&nextPageToken)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !bytes.Equal(nonce, nextPageToken.Nonce) {
		return 0, ErrChangedParameters
	}

	return nextPageToken.Page, nil
}

// Encode takes a page number for the current response, the next page size, and
// a nonce, and returns an encoded `next_page_token` UTF8 string to pass to a
// REST client as a way to continue a list query at the next page after this
// one.  If the encoding fails, an error is returned instead and the token is
// undefined.  The nonce must be constant for the same query parameters, e.g. a
// hash of the filter expression string.  The page size may change between
// requests.
func Encode(page int, nonce []byte) (next_page_token string, err error) {
	token := nextPageToken{
		Page:  page + 1,
		Nonce: nonce,
	}
	var b bytes.Buffer
	gobEnc := gob.NewEncoder(&b)
	err = gobEnc.Encode(token)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b.Bytes()), nil
}
