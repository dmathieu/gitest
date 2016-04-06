package gitest

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func (s *Server) refsEndpoint(w http.ResponseWriter, r *http.Request) {
	service := r.FormValue("service")
	if !validService(service) {
		w.WriteHeader(400)
	}

	w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))

	w.Write(packetWrite(fmt.Sprintf("# service=%s\n", service)))

	w.Write([]byte("0000"))

	command := exec.Command(service, "--stateless-rpc", "--advertise-refs", s.template.Folder)
	command.Stdout = w
	command.Stderr = w
	command.Run()
}

func (s *Server) serviceEndpoint(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get(":service")
	if !validService(service) {
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", service))

	defer r.Body.Close()
	command := exec.Command(service, "--stateless-rpc", s.template.Folder)
	command.Stdin = r.Body
	command.Stdout = w
	command.Stderr = w
	command.Run()
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)
	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}
	return []byte(s + str)
}

func validService(name string) bool {
	return name == "git-upload-pack" || name == "git-receive-pack"
}
