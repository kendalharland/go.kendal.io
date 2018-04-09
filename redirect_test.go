package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	main "go.kendal.io/go.kendal.io"
)

func TestRedirect(t *testing.T) {
	t.Run("Browser", func(t *testing.T) {
		expectRedirect(t,
			http.MethodGet,
			"https://go.kendal.io",
			"https://godoc.org/?q=go.kendal.io",
			http.StatusTemporaryRedirect)
		expectRedirect(t,
			http.MethodGet,
			"https://go.kendal.io/",
			"https://godoc.org/?q=go.kendal.io",
			http.StatusTemporaryRedirect)
		expectRedirect(t,
			http.MethodGet,
			"https://go.kendal.io?go-get=0",
			"https://godoc.org/?q=go.kendal.io",
			http.StatusTemporaryRedirect)
		expectRedirect(t,
			http.MethodGet,
			"https://go.kendal.io/foo",
			"https://godoc.org/go.kendal.io/foo",
			http.StatusTemporaryRedirect)
	})

	t.Run("HTTP", func(t *testing.T) {
		expectRedirect(t,
			http.MethodGet,
			"http://go.kendal.io/foo?go-get=1",
			"https://go.kendal.io/foo?go-get=1",
			http.StatusMovedPermanently)
		expectRedirect(t,
			http.MethodGet,
			"http://go.kendal.io",
			"https://go.kendal.io",
			http.StatusMovedPermanently)
		expectRedirect(t,
			http.MethodGet,
			"http://go.kendal.io/",
			"https://go.kendal.io/",
			http.StatusMovedPermanently)
	})

	t.Run("Method", func(t *testing.T) {
		expectRedirect(t,
			http.MethodHead,
			"https://go.kendal.io",
			"",
			http.StatusMethodNotAllowed)
		expectRedirect(t,
			http.MethodPost,
			"https://go.kendal.io",
			"",
			http.StatusMethodNotAllowed)
		expectRedirect(t,
			http.MethodPut,
			"https://go.kendal.io",
			"",
			http.StatusMethodNotAllowed)
		expectRedirect(t,
			http.MethodPatch,
			"https://go.kendal.io",
			"",
			http.StatusMethodNotAllowed)
		expectRedirect(t,
			http.MethodDelete,
			"https://go.kendal.io",
			"",
			http.StatusMethodNotAllowed)
		expectRedirect(t,
			http.MethodConnect,
			"https://go.kendal.io",
			"",
			http.StatusMethodNotAllowed)
		expectRedirect(t,
			http.MethodOptions,
			"https://go.kendal.io",
			"",
			http.StatusMethodNotAllowed)
		expectRedirect(t,
			http.MethodTrace,
			"https://go.kendal.io",
			"",
			http.StatusMethodNotAllowed)
	})

	t.Run("GitHub", func(t *testing.T) {
		expectRedirectMetaTag(t,
			"https://go.kendal.io/foo?go-get=1",
			"<meta name=\"go-import\" content=\"go.kendal.io/foo git https://github.com/kharland/foo\"/>",
			http.StatusOK)
		// extra slash after "/foo"
		expectRedirectMetaTag(t,
			"https://go.kendal.io/foo/?go-get=1",
			"<meta name=\"go-import\" content=\"go.kendal.io/foo git https://github.com/kharland/foo\"/>",
			http.StatusOK)
		expectRedirectMetaTag(t,
			"https://go.kendal.io/foo/cmd/bar?go-get=1",
			"<meta name=\"go-import\" content=\"go.kendal.io/foo git https://github.com/kharland/foo\"/>",
			http.StatusOK)
	})
}

func expectRedirect(t *testing.T, method string, from string, to string, code int) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest(method, from, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	main.Redirect(res, req)
	if res.Code != code {
		t.Errorf("expected status %v but got %v", code, res.Code)
	}

	expectedLoc := to
	actualLoc := res.HeaderMap.Get("Location")
	if actualLoc != expectedLoc {
		t.Errorf("expected redirect to %s.  Got %s", expectedLoc, actualLoc)
	}
}

func expectRedirectMetaTag(t *testing.T, from string, tag string, code int) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, from, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	main.Redirect(res, req)
	if res.Code != code {
		t.Errorf("expected status %v but got %v", code, res.Code)
	}

	body := res.Body.String()
	if !strings.Contains(body, tag) {
		t.Errorf("expected body to have tag %s. Got %s", tag, body)
	}
}
