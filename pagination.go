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
	"reflect"
)

const MAX_TOKEN_LENGTH = 1024

var ErrChangedParameters = errors.New("parameters changed between method calls invalidating this token")
var ErrInvalidToken = errors.New("invalid token")
var ErrPagination = errors.New("unable to compute next_page_token")

type nextPageToken[Cursor any] struct {
	Cursor Cursor
	Nonce  []byte
}

func zeroOf[Cursor any]() Cursor {
	return reflect.Zero(reflect.TypeFor[Cursor]()).Interface().(Cursor)
}

// Decode takes a `next_page_token` UTF8 string and a nonce and returns a
// decoded Cursor, or an error explaining why this token was not
// valid.  The nonce must be constant for the same query parameters, e.g. a
// hash of a filter expression string.  If the token is empty, the Cursor is
// zero.
func Decode[Cursor any](token string, nonce []byte) (Cursor, error) {
	if token == "" {
		return zeroOf[Cursor](), nil
	}
	if len(token) > MAX_TOKEN_LENGTH {
		return zeroOf[Cursor](), fmt.Errorf("%w: token too long", ErrInvalidToken)
	}
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return zeroOf[Cursor](), fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	bs := bytes.NewBuffer(b)

	var nextPageToken nextPageToken[Cursor]
	gDec := gob.NewDecoder(bs)
	err = gDec.Decode(&nextPageToken)
	if err != nil {
		return zeroOf[Cursor](), fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !bytes.Equal(nonce, nextPageToken.Nonce) {
		return zeroOf[Cursor](), ErrChangedParameters
	}

	return nextPageToken.Cursor, nil
}

// Encode takes a cursor and
// a nonce, and returns an encoded `next_page_token` UTF8 string to pass to a
// REST client as a way to continue a list query at the next page after this
// one.  If the encoding fails, an error is returned instead and the token is
// undefined.  The nonce must be constant for the same query parameters, e.g. a
// hash of the filter expression string.
//
// See https://google.aip.dev/158 for an explanation.
func Encode[Cursor any](lastCursor Cursor, nonce []byte) (next_page_token string, err error) {
	token := nextPageToken[Cursor]{
		Cursor: lastCursor,
		Nonce:  nonce,
	}

	var b bytes.Buffer
	gobEnc := gob.NewEncoder(&b)
	err = gobEnc.Encode(token)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrPagination, err)
	}
	return base64.RawURLEncoding.EncodeToString(b.Bytes()), nil
}
