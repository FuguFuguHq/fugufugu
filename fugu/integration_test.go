package fugu

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
<link rel="stylesheet" href="main.css">
</head>
<body>
<h1>Hello World</h1>
<img src="https://www.amazingcto.com/AmazingCTO_Logo_White.svg">
<p class="description">This is a test page</p>
<p class="description">This is a test paragraph</p>
</body>
</html>
		`))
	})
	mux.HandleFunc("/main.css", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(`
@import url('https://fonts.googleapis.com/css?family=Muli')
body {
}
`))
	})
	return httptest.NewServer(mux)
}

func TestCollectorVisit(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	externals := make(map[string]Privacy)

	c := NewCollector(ts.URL, externals, false)
	err := c.Visit(ts.URL)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 2, len(externals))
	assert.Equal(t, "Image", externals["https://www.amazingcto.com/AmazingCTO_Logo_White.svg"].Typ)
	assert.Equal(t, "Css", externals["https://fonts.googleapis.com/css?family=Muli"].Typ)
}
