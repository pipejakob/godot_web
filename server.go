package godot_web

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"path"
)

type Server struct {
	Dir                  string
	Port                 int
	AllowExternalTraffic bool
	TLSCertFile          string
	TLSKeyFile           string
}

func New(dir string, port int, allowExternalTraffic bool, tlsCertFile, tlsKeyFile string) *Server {
	return &Server{
		Dir:                  dir,
		Port:                 port,
		AllowExternalTraffic: allowExternalTraffic,
		TLSCertFile:          tlsCertFile,
		TLSKeyFile:           tlsKeyFile,
	}
}

func (s Server) Run() error {
	httpServer := &http.Server{
		Addr:    s.listenAddress(),
		Handler: s,
	}

	link, err := s.link()
	if err != nil {
		return fmt.Errorf("error generating server link: %v", err)
	}

	fmt.Printf("Server link: %s\n", link)

	if !s.AllowExternalTraffic {
		return httpServer.ListenAndServe()
	}

	if s.TLSCertFile != "" && s.TLSKeyFile != "" {
		return httpServer.ListenAndServeTLS(s.TLSCertFile, s.TLSKeyFile)
	}

	ip, err := s.ip()
	if err != nil {
		return fmt.Errorf("error getting IP address: %v", err)
	}

	cert, err := generateSelfSignedCertificate(ip)
	if err != nil {
		return fmt.Errorf("error generating certificate: %v", err)
	}

	httpServer.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},
	}

	return httpServer.ListenAndServeTLS("", "")
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullPath := path.Join(s.Dir, r.URL.Path)

	w.Header().Add("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Add("Cross-Origin-Embedder-Policy", "require-corp")

	http.ServeFile(w, r, fullPath)
}

func (s Server) listenAddress() string {
	if s.AllowExternalTraffic {
		return fmt.Sprintf(":%d", s.Port)
	} else {
		return fmt.Sprintf("127.0.0.1:%d", s.Port)
	}
}

func (s Server) link() (string, error) {
	if s.AllowExternalTraffic {
		ip, err := s.ip()
		if err != nil {
			return "", fmt.Errorf("error getting IP address: %v", err)
		}

		return fmt.Sprintf("https://%s:%d/", ip, s.Port), nil
	} else {
		return fmt.Sprintf("http://127.0.0.1:%d/", s.Port), nil
	}
}

func (s Server) ip() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("error dialing: %v", err)
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}
