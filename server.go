package godot_web

import (
	"fmt"
	"net/http"
	"path"
)

type Server struct {
	Dir  string
	Port int
}

func New(dir string, port int) *Server {
	return &Server{
		Dir:  dir,
		Port: port,
	}
}

func (s Server) Run() error {
	http.HandleFunc("/", s.serveFile)

	link, err := s.link()
	if err != nil {
		return err
	}

	fmt.Printf("Server link: %s\n", link)

	return http.ListenAndServe(s.listenAddress(), nil)
}

func (s Server) serveFile(w http.ResponseWriter, r *http.Request) {
	fullPath := path.Join(s.Dir, r.URL.Path)

	w.Header().Add("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Add("Cross-Origin-Embedder-Policy", "require-corp")

	http.ServeFile(w, r, fullPath)
}

func (s Server) listenAddress() string {
	return fmt.Sprintf("127.0.0.1:%d", s.Port)
}

func (s Server) link() (string, error) {
	return fmt.Sprintf("http://127.0.0.1:%d/", s.Port), nil
}
