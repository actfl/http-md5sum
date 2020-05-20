package httpsum

import (
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type siteResponse struct {
	site    string
	success bool
	err     string
	md5     string
}
