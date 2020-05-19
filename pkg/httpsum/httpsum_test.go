package httpsum

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	// just in case you want default correct return value
	return &http.Response{}, nil
}

func TestHttpSum_get(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		timeout    int

		want        [md5.Size]byte
		wantErr     bool
		expectedErr error
	}{
		{
			name:       "0: successful result",
			statusCode: http.StatusOK,
			body:       "hello world",
			want:       md5.Sum([]byte("hello world")),
		},
		{
			name:        "1: url not found",
			statusCode:  http.StatusNotFound,
			wantErr:     true,
			expectedErr: httpStatusError,
		},
		{
			name:        "2: timeout",
			timeout:     2,
			wantErr:     true,
			expectedErr: timeoutError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if tt.timeout > 0 {
						time.Sleep(time.Duration(tt.timeout) * time.Second)
					}
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(tt.body))),
					}, nil
				},
			}

			c := Config{
				Client:   &client,
				Parallel: 1,
				Timeout:  1,
			}

			h, err := New(c)
			if err != nil {
				t.Errorf("expect nil, got %v", err)
			}

			got, err := h.get("www.example.com")
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("get() error = %v, expected error type %v", err, tt.expectedErr)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
