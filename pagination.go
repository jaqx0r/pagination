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
	Offset  int
	Nonce []byte
}

// Decode takes a `next_page_token` UTF8 string and a nonce and returns a
// decoded offset number, or an error explaining why this token was not
// valid.  The nonce must be constant for the same query parameters, e.g. a
// hash of a filter expression string.  If the token is empty, the offset is
// zero.  
func Decode(token string, nonce []byte) (offset int, err error) {
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

	return nextPageToken.Offset, nil
}

// Encode takes the last offset, the page size for the current response, and
// a nonce, and returns an encoded `next_page_token` UTF8 string to pass to a
// REST client as a way to continue a list query at the next page after this
// one.  If the encoding fails, an error is returned instead and the token is
// undefined.  The nonce must be constant for the same query parameters, e.g. a
// hash of the filter expression string.
func Encode(lastOffset, pageSize int, nonce []byte) (next_page_token string, err error) {
	token := nextPageToken{
		Offset: lastOffset + pageSize,
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
