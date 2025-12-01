# pagination

A library to assist with REST API pagination, implemented in Go.

An [AIP-158](https://google.aip.dev/158) compliant pagination function needs to transform an incoming `size`, `token`, and request signature into an `offset` and `limit` suitable for continuing a query against storage.

Likewise the function needs to compute a new `next_page_token` for the response, such that the client can return it in a followup request.

The `next_page_token` is not required, which means the pagination starts at offset 0.

The `page_size` defines the limit of responses and defaults to a caller-provided value.

A response may contain fewer than `page_size` resources in the collection, possibly even zero, but only when `next_page_token` is empty can the client assume there are no more resources in the collection.

The `next_page_token` must be an opaque string.

The request signature (i.e. any other query parameters) should be constant across all queries using this sequence of page tokens, so the function should return an error if the signature changes for a given token.
