# pagination

A library to assist with REST API pagination, implemented in Go.

## API

```go
pagination.Decode(token string, nonce []byte) (offset int, err error)
```

Decode a page `token` provided in the request parameters by the caller, to get the cursor offset.  Errors indicate invalid tokens and should return an Invalid Argument HTTP response code.

The `nonce` is a signature of the request query parameters, for example a hash, that the function will use to determine if the offset is still valid with respect to those query parameters.

```go
pagination.Encode(offset, pageSize int, nonce []byte) (token string, err error)
```

Encode the last page offset (from a prior `Decode()`), the current page size, and the previously computed nonce into a token string to return to the caller in the response payload, so that they may return it the next time.

## Specification

An [AIP-158](https://google.aip.dev/158) compliant pagination function needs to transform an incoming `token` and request signature into a page `offset` suitable for continuing a query against storage.

Likewise the function needs to compute a new `next_page_token` for the response, such that the client can return it in a followup request.

The `next_page_token` is not required, which means the pagination starts at offset 0.

A response may contain fewer than `page_size` resources in the collection, possibly even zero, but only when `next_page_token` is empty can the client assume there are no more resources in the collection.

The `next_page_token` must be an opaque string.

The request signature (i.e. any other query parameters) should be constant across all queries using this sequence of page tokens, so the function should return an error if the signature changes for a given token.

## Example

https://github.com/jaqx0r/pagination/blob/51750129674167c0d34c594467b107a9f66b1aa0/pagination_test.go#L101-L117
