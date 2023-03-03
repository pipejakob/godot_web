package godot_web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

func TestServer_serveFile_ReturnsOKStatus(t *testing.T) {
	dir := t.TempDir()
	createFile(t, dir, "test.html", []byte("hello world"))
	r, w := createRequestResponse("/test.html")
	server := New(dir, 8000)

	server.serveFile(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("serveFile() wanted response code %d, got: %d", expected, w.Code)
	}
}

func TestServer_serveFile_NestedFileReturnsOKStatus(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(path.Join(dir, "my/long/path/to"), 0700)
	createFile(t, dir, "my/long/path/to/test.html", []byte("hello world"))
	r, w := createRequestResponse("/my/long/path/to/test.html")
	server := New(dir, 8000)

	server.serveFile(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("serveFile() wanted response code %d, got: %d", expected, w.Code)
	}
}

func TestServer_serveFile_IndexHTMLReturnsRedirect(t *testing.T) {
	dir := t.TempDir()
	createFile(t, dir, "index.html", []byte("hello world"))
	r, w := createRequestResponse("http://127.0.0.1/index.html")
	server := New(dir, 8000)

	server.serveFile(w, r)

	expectedCode := http.StatusMovedPermanently
	expectedLocation := "./"
	if w.Code != expectedCode {
		t.Errorf("serveFile() wanted response code %d, got: %d", expectedCode, w.Code)
	}
	if w.HeaderMap.Get("Location") != expectedLocation {
		t.Errorf("serveFile() wanted response header Location = %q, got: %q",
			expectedLocation, w.HeaderMap.Get("Location"))
	}
}

func TestServer_serveFile_NestedIndexHTMLReturnsRedirect(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(path.Join(dir, "my/project/dir"), 0700)
	createFile(t, dir, "my/project/dir/index.html", []byte("hello world"))
	r, w := createRequestResponse("http://127.0.0.1/my/project/dir/index.html")
	server := New(dir, 8000)

	server.serveFile(w, r)

	expectedCode := http.StatusMovedPermanently
	expectedLocation := "./"
	if w.Code != expectedCode {
		t.Errorf("serveFile() wanted response code %d, got: %d", expectedCode, w.Code)
	}
	if w.HeaderMap.Get("Location") != expectedLocation {
		t.Errorf("serveFile() wanted response header Location = %q, got: %q",
			expectedLocation, w.HeaderMap.Get("Location"))
	}
}

func TestServer_serveFile_DirectoryReturnsFileListing(t *testing.T) {
	dir := t.TempDir()
	contents := []byte("hello world")
	createFile(t, dir, "test.html", contents)
	createFile(t, dir, "game.js", contents)
	createFile(t, dir, "README.txt", contents)
	r, w := createRequestResponse("http://127.0.0.1/")
	server := New(dir, 8000)

	server.serveFile(w, r)

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("serveFile() wanted response code %d, got: %d", expectedCode, w.Code)
	}
	body := string(w.Body.Bytes())
	if !strings.Contains(body, "<a href=\"test.html\">") {
		t.Errorf("serveFile() could not find response body link for \"test.html\"")
	}
	if !strings.Contains(body, "<a href=\"game.js\">") {
		t.Errorf("serveFile() could not find response body link for \"game.js\"")
	}
	if !strings.Contains(body, "<a href=\"README.txt\">") {
		t.Errorf("serveFile() could not find response body link for \"README.txt\"")
	}
}

func TestServer_serveFile_DirectoryWithIndexHTMLReturnsIt(t *testing.T) {
	dir := t.TempDir()
	contents := []byte("hello world")
	createFile(t, dir, "index.html", contents)
	r, w := createRequestResponse("http://127.0.0.1/")
	server := New(dir, 8000)

	server.serveFile(w, r)

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("serveFile() wanted response code %d, got: %d", expectedCode, w.Code)
	}
	if !bytes.Equal(w.Body.Bytes(), contents) {
		t.Errorf("serveFile() wanted body %v, got: %v", contents, w.Body.Bytes())
	}
}

func TestServer_serveFile_URLWithDotDotReturnsBadRequest(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(path.Join(dir, "my/project/dir"), 0700)
	createFile(t, dir, "my/project/dir/test.html", []byte("hello world"))
	r, w := createRequestResponse("http://127.0.0.1/my/project/../project/dir/test.html")
	server := New(dir, 8000)

	server.serveFile(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("serveFile() wanted response code %d, got: %d", expected, w.Code)
	}
}

func TestServer_serveFile_MissingFileReturnsNotFoundStatus(t *testing.T) {
	dir := t.TempDir()
	r, w := createRequestResponse("/missing.html")
	server := New(dir, 8000)

	server.serveFile(w, r)

	expected := http.StatusNotFound
	if w.Code != expected {
		t.Errorf("serveFile() wanted response code %d, got: %d", expected, w.Code)
	}
}

func TestServer_serveFile_FileReturnsCorrectBody(t *testing.T) {
	contents := []byte("hello world")
	dir := t.TempDir()
	createFile(t, dir, "test.html", contents)
	r, w := createRequestResponse("/test.html")
	server := New(dir, 8000)

	server.serveFile(w, r)

	if !bytes.Equal(w.Body.Bytes(), contents) {
		t.Errorf("serveFile() wanted body %v, got: %v", contents, w.Body.Bytes())
	}
}

func TestServer_serveFile_FileReturnsCrossOriginHeaders(t *testing.T) {
	dir := t.TempDir()
	createFile(t, dir, "test.html", []byte("hello world"))
	r, w := createRequestResponse("/test.html")
	server := New(dir, 8000)

	server.serveFile(w, r)
	openerPolicy := w.HeaderMap.Values("Cross-Origin-Opener-Policy")
	embedderPolicy := w.HeaderMap.Values("Cross-Origin-Embedder-Policy")

	expectedOwnerPolicy := "same-origin"
	expectedEmbedderPolicy := "require-corp"
	if len(openerPolicy) != 1 {
		t.Errorf("serveFile() expected 1 Cross-Origin-Opener-Policy header, got %d", len(openerPolicy))
	}
	if len(openerPolicy) >= 1 && openerPolicy[0] != expectedOwnerPolicy {
		t.Errorf("serveFile() expected Cross-Origin-Opener-Policy header value %q, got %q",
			expectedOwnerPolicy, openerPolicy[0])
	}
	if len(embedderPolicy) != 1 {
		t.Errorf("serveFile() expected 1 Cross-Origin-Embedder-Policy header, got %d", len(embedderPolicy))
	}
	if len(embedderPolicy) >= 1 && embedderPolicy[0] != expectedEmbedderPolicy {
		t.Errorf("serveFile() expected Cross-Origin-Embedder-Policy header value %q, got %q",
			expectedEmbedderPolicy, embedderPolicy[0])
	}
}

func TestServer_listenAddress_IsCorrect(t *testing.T) {
	dir := t.TempDir()
	server := New(dir, 8000)

	addr := server.listenAddress()

	expected := "127.0.0.1:8000"
	if addr != expected {
		t.Errorf("listenAddress() wanted %q, got %q", expected, addr)
	}
}

func createFile(t *testing.T, dir string, name string, contents []byte) {
	fullPath := path.Join(dir, name)

	if err := os.WriteFile(fullPath, []byte(contents), 0600); err != nil {
		t.Fatalf("error creating test file: %v", err)
	}
}

func createRequestResponse(path string) (*http.Request, *httptest.ResponseRecorder) {
	r := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()

	return r, w
}
