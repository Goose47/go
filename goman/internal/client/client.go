package client

import (
	"io"
	"net/http"
	"strings"
	"sync"
)

const BASE_URL = "https://pkg.go.dev"

var cache map[string]string
var mu sync.Mutex

func CachedFetch(url string) (io.Reader, error) {
	url = BASE_URL + url
	mu.Lock()
	if cache == nil {
		cache = make(map[string]string)
	}
	cached, hit := cache[url]
	mu.Unlock()

	if hit {
		return strings.NewReader(cached), nil
	}

	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(result.Body)
	result.Body.Close()

	bytesStr := string(bytes)

	mu.Lock()
	cache[url] = bytesStr
	mu.Unlock()

	return strings.NewReader(bytesStr), nil
}
