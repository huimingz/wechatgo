package mock

import "net/http"

func NewMockedHTTP() *http.Client {
	return http.DefaultClient
}
