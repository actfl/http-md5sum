package httpsum

import (
	"crypto/md5"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type siteResponse struct {
	site    string
	success bool
	md5     [md5.Size]byte
	err     error
}
