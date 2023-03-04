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

func TestServer_ServeHTTP_ReturnsOKStatus(t *testing.T) {
	dir := t.TempDir()
	createFile(t, dir, "test.html", []byte("hello world"))
	r, w := createRequestResponse("/test.html")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("ServeHTTP() wanted response code %d, got: %d", expected, w.Code)
	}
}

func TestServer_ServeHTTP_NestedFileReturnsOKStatus(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(path.Join(dir, "my/long/path/to"), 0700)
	createFile(t, dir, "my/long/path/to/test.html", []byte("hello world"))
	r, w := createRequestResponse("/my/long/path/to/test.html")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("ServeHTTP() wanted response code %d, got: %d", expected, w.Code)
	}
}

func TestServer_ServeHTTP_IndexHTMLReturnsRedirect(t *testing.T) {
	dir := t.TempDir()
	createFile(t, dir, "index.html", []byte("hello world"))
	r, w := createRequestResponse("http://127.0.0.1/index.html")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	expectedCode := http.StatusMovedPermanently
	expectedLocation := "./"
	if w.Code != expectedCode {
		t.Errorf("ServeHTTP() wanted response code %d, got: %d", expectedCode, w.Code)
	}
	if w.HeaderMap.Get("Location") != expectedLocation {
		t.Errorf("ServeHTTP() wanted response header Location = %q, got: %q",
			expectedLocation, w.HeaderMap.Get("Location"))
	}
}

func TestServer_ServeHTTP_NestedIndexHTMLReturnsRedirect(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(path.Join(dir, "my/project/dir"), 0700)
	createFile(t, dir, "my/project/dir/index.html", []byte("hello world"))
	r, w := createRequestResponse("http://127.0.0.1/my/project/dir/index.html")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	expectedCode := http.StatusMovedPermanently
	expectedLocation := "./"
	if w.Code != expectedCode {
		t.Errorf("ServeHTTP() wanted response code %d, got: %d", expectedCode, w.Code)
	}
	if w.HeaderMap.Get("Location") != expectedLocation {
		t.Errorf("ServeHTTP() wanted response header Location = %q, got: %q",
			expectedLocation, w.HeaderMap.Get("Location"))
	}
}

func TestServer_ServeHTTP_DirectoryReturnsFileListing(t *testing.T) {
	dir := t.TempDir()
	contents := []byte("hello world")
	createFile(t, dir, "test.html", contents)
	createFile(t, dir, "game.js", contents)
	createFile(t, dir, "README.txt", contents)
	r, w := createRequestResponse("http://127.0.0.1/")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("ServeHTTP() wanted response code %d, got: %d", expectedCode, w.Code)
	}
	body := string(w.Body.Bytes())
	if !strings.Contains(body, "<a href=\"test.html\">") {
		t.Errorf("ServeHTTP() could not find response body link for \"test.html\"")
	}
	if !strings.Contains(body, "<a href=\"game.js\">") {
		t.Errorf("ServeHTTP() could not find response body link for \"game.js\"")
	}
	if !strings.Contains(body, "<a href=\"README.txt\">") {
		t.Errorf("ServeHTTP() could not find response body link for \"README.txt\"")
	}
}

func TestServer_ServeHTTP_DirectoryWithIndexHTMLReturnsIt(t *testing.T) {
	dir := t.TempDir()
	contents := []byte("hello world")
	createFile(t, dir, "index.html", contents)
	r, w := createRequestResponse("http://127.0.0.1/")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("ServeHTTP() wanted response code %d, got: %d", expectedCode, w.Code)
	}
	if !bytes.Equal(w.Body.Bytes(), contents) {
		t.Errorf("ServeHTTP() wanted body %v, got: %v", contents, w.Body.Bytes())
	}
}

func TestServer_ServeHTTP_URLWithDotDotReturnsBadRequest(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(path.Join(dir, "my/project/dir"), 0700)
	createFile(t, dir, "my/project/dir/test.html", []byte("hello world"))
	r, w := createRequestResponse("http://127.0.0.1/my/project/../project/dir/test.html")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("ServeHTTP() wanted response code %d, got: %d", expected, w.Code)
	}
}

func TestServer_ServeHTTP_MissingFileReturnsNotFoundStatus(t *testing.T) {
	dir := t.TempDir()
	r, w := createRequestResponse("/missing.html")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	expected := http.StatusNotFound
	if w.Code != expected {
		t.Errorf("ServeHTTP() wanted response code %d, got: %d", expected, w.Code)
	}
}

func TestServer_ServeHTTP_FileReturnsCorrectBody(t *testing.T) {
	contents := []byte("hello world")
	dir := t.TempDir()
	createFile(t, dir, "test.html", contents)
	r, w := createRequestResponse("/test.html")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)

	if !bytes.Equal(w.Body.Bytes(), contents) {
		t.Errorf("ServeHTTP() wanted body %v, got: %v", contents, w.Body.Bytes())
	}
}

func TestServer_ServeHTTP_FileReturnsCrossOriginHeaders(t *testing.T) {
	dir := t.TempDir()
	createFile(t, dir, "test.html", []byte("hello world"))
	r, w := createRequestResponse("/test.html")
	server := newDefaultServer(dir)

	server.ServeHTTP(w, r)
	openerPolicy := w.HeaderMap.Values("Cross-Origin-Opener-Policy")
	embedderPolicy := w.HeaderMap.Values("Cross-Origin-Embedder-Policy")

	expectedOwnerPolicy := "same-origin"
	expectedEmbedderPolicy := "require-corp"
	if len(openerPolicy) != 1 {
		t.Errorf("ServeHTTP() expected 1 Cross-Origin-Opener-Policy header, got %d", len(openerPolicy))
	}
	if len(openerPolicy) >= 1 && openerPolicy[0] != expectedOwnerPolicy {
		t.Errorf("ServeHTTP() expected Cross-Origin-Opener-Policy header value %q, got %q",
			expectedOwnerPolicy, openerPolicy[0])
	}
	if len(embedderPolicy) != 1 {
		t.Errorf("ServeHTTP() expected 1 Cross-Origin-Embedder-Policy header, got %d", len(embedderPolicy))
	}
	if len(embedderPolicy) >= 1 && embedderPolicy[0] != expectedEmbedderPolicy {
		t.Errorf("ServeHTTP() expected Cross-Origin-Embedder-Policy header value %q, got %q",
			expectedEmbedderPolicy, embedderPolicy[0])
	}
}

func TestServer_listenAddress_IsCorrectForLocalhost(t *testing.T) {
	dir := t.TempDir()
	server := New(dir, 8000, false, "", "")

	addr := server.listenAddress()

	expected := "127.0.0.1:8000"
	if addr != expected {
		t.Errorf("listenAddress() wanted %q, got %q", expected, addr)
	}
}

func TestServer_listenAddress_IsCorrectForExternal(t *testing.T) {
	dir := t.TempDir()
	server := New(dir, 11075, true, "", "")

	addr := server.listenAddress()

	expected := ":11075"
	if addr != expected {
		t.Errorf("listenAddress() wanted %q, got %q", expected, addr)
	}
}

func TestServer_link_InternalUsesHTTP(t *testing.T) {
	dir := t.TempDir()
	server := New(dir, 8080, false, "", "")

	link, err := server.link()

	if err != nil {
		t.Errorf("link() returned unwanted error: %v", err)
	}
	expectedStart := "http://"
	if !strings.HasPrefix(link, expectedStart) {
		t.Errorf("link() wanted prefix %q, got %q", expectedStart, link)
	}
}

func TestServer_link_ExternalUsesHTTPS(t *testing.T) {
	dir := t.TempDir()
	server := New(dir, 8443, true, "", "")

	link, err := server.link()

	if err != nil {
		t.Errorf("link() returned unwanted error: %v", err)
	}
	expectedStart := "https://"
	if !strings.HasPrefix(link, expectedStart) {
		t.Errorf("link() wanted prefix %q, got %q", expectedStart, link)
	}
}

func newDefaultServer(dir string) *Server {
	return New(dir, 8000, false, "", "")
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
