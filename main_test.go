package main

import (
	"fmt"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestRewriteHandler(t *testing.T) {
	const checksum = "test-checksum"
	const port = 12345
	testdata := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte(`Sample File with http://127.0.0.1:6120/foo and wss://127.0.0.1:6120/bar`),
		},
		"other.html": &fstest.MapFile{
			Data: []byte(`Another File with http://127.0.0.1:6120/foo and wss://127.0.0.1:6120/bar`),
		},
		"path/to/javascript.js": &fstest.MapFile{
			Data: []byte(`/* Hello http://127.0.0.1:6120/foo and wss://127.0.0.1:6120/bar */`),
		},
		"image.png": &fstest.MapFile{
			Data: []byte(`Not modified http://127.0.0.1:6120/image.png`),
		},
	}

	tests := []struct {
		// name is the name of the test case, used for reporting.
		name string
		// path is the URL path to request.
		path string
		// header is the value for the If-None-Match header; if empty, the header is not sent.
		header string
		// expectNotModified indicates whether the handler should respond with 304 Not Modified.
		expectNotModified bool
		// expected is the expected body of the response.
		expected string
		// contentType is the expected Content-Type header of the response; if empty, it is not checked.
		contentType string
	}{
		{
			name:        "index.html",
			path:        "/index.html",
			expected:    `Sample File with http://127.0.0.1:12345/foo and wss://127.0.0.1:12345/bar`,
			contentType: mime.TypeByExtension(".html"),
		},
		{
			name:     "root path",
			path:     "/",
			expected: `Sample File with http://127.0.0.1:12345/foo and wss://127.0.0.1:12345/bar`,
		},
		{
			name:        "non-existent file",
			path:        "/this/path/does/not/exist.png",
			expected:    `Sample File with http://127.0.0.1:12345/foo and wss://127.0.0.1:12345/bar`,
			contentType: mime.TypeByExtension(".html"),
		},
		{
			name:     "other html file",
			path:     "/other.html",
			expected: `Another File with http://127.0.0.1:12345/foo and wss://127.0.0.1:12345/bar`,
		},
		{
			name:        "javascript file",
			path:        "/path/to/javascript.js",
			expected:    `/* Hello http://127.0.0.1:12345/foo and wss://127.0.0.1:12345/bar */`,
			contentType: mime.TypeByExtension(".js"),
		},
		{
			name:        "image file",
			path:        "/image.png",
			expected:    "Not modified http://127.0.0.1:6120/image.png",
			contentType: mime.TypeByExtension(".png"),
		},
		{
			name:             "not modified",
			path:             "/index.html",
			header:           fmt.Sprintf(`W/"%s"`, checksum),
			expectNotModified: true,
		},
		{
			name:             "not modified on rewrite path",
			path:             "/c/local/explorer",
			header:           fmt.Sprintf(`W/"%s"`, checksum),
			expectNotModified: true,
		},
		{
			// We should allow strong ETags even though we send a weak one.
			name:             "strong not modified",
			path:             "/image.png",
			header:           fmt.Sprintf(`"%s"`, checksum),
			expectNotModified: true,
		},
		{
			// We should allow strong ETags even though we send a weak one.
			name:             "strong not modified on rewrite path",
			path:             "/c/local/explorer",
			header:           fmt.Sprintf(`"%s"`, checksum),
			expectNotModified: true,
		},
	}

	handler, err := rewriteHandler(port, testdata, checksum)
	if err != nil {
		t.Fatalf("rewriteHandler returned error: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to the handler.
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, tt.path, http.NoBody)
			if tt.header != "" {
				req.Header.Set("If-None-Match", tt.header)
			}

			resp := httptest.NewRecorder()
			handler(resp, req)

			if tt.contentType != "" {
				if resp.Header().Get("Content-Type") != tt.contentType {
					t.Errorf("unexpected Content-Type: got %q, want %q", resp.Header().Get("Content-Type"), tt.contentType)
				}
			}

			if tt.expectNotModified {
				if resp.Code != http.StatusNotModified {
					t.Errorf("unexpected status code: got %d, want %d", resp.Code, http.StatusNotModified)
				}
				if resp.Body.Len() != 0 {
					t.Errorf("unexpected body: got %q, want empty", resp.Body.String())
				}
				return
			}

			if resp.Code != http.StatusOK {
				t.Errorf("unexpected status code: got %d, want %d", resp.Code, http.StatusOK)
			}
			if resp.Body.String() != tt.expected {
				t.Errorf("unexpected body: got %q, want %q", resp.Body.String(), tt.expected)
			}
		})
	}
}

func TestIsNotModified(t *testing.T) {
	// The checksum is normally set at build time; since we do not have that at
	// test time, just set it to some fixed value that does not occur in tests.
	const desiredTag = "test-checksum"
	tests := []struct {
		name        string
		ifNoneMatch string
		expected    bool
	}{
		{
			name:        "empty header",
			ifNoneMatch: "",
			expected:    false,
		},
		{
			name:        "matching tag",
			ifNoneMatch: fmt.Sprintf(`W/"%s"`, desiredTag),
			expected:    true,
		},
		{
			name:        "matching tag with whitespace",
			ifNoneMatch: fmt.Sprintf(`  W/"%s"  `, desiredTag),
			expected:    true,
		},
		{
			name:        "misplaced whitespace",
			ifNoneMatch: fmt.Sprintf(`W/  "%s"  `, desiredTag),
			expected:    false,
		},
		{
			name:        "multiple tags with match",
			ifNoneMatch: fmt.Sprintf(`foo, W/"%s", bar`, desiredTag),
			expected:    true,
		},
		{
			name:        "multiple tags without match",
			ifNoneMatch: `foo, bar, baz`,
			expected:    false,
		},
		{
			name:        "wildcard match",
			ifNoneMatch: `*`,
			expected:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNotModified(tt.ifNoneMatch, desiredTag)
			if result != tt.expected {
				t.Errorf("isNotModified(%q) = %v; want %v", tt.ifNoneMatch, result, tt.expected)
			}
		})
	}
}
