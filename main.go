package gitest

import (
	"math/rand"
	"net/http/httptest"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Server is an HTTP listener able to respond to git endpoints
type Server struct {
	*httptest.Server
	template *template

	ValidRepo      string
	NotAllowedRepo string
}

// NewServer creates a new Server object
func NewServer(template string) (*Server, error) {
	s, err := NewUnstartedServer(template)
	if err != nil {
		return nil, err
	}
	s.Server.Start()

	return s, nil
}

// NewUnstartedServer creates a new server but doesn't start it
func NewUnstartedServer(template string) (*Server, error) {
	rand.Seed(time.Now().UnixNano())

	t, err := newTemplate(template)
	if err != nil {
		return nil, err
	}

	s := &Server{
		template:  t,
		ValidRepo: generateRepoName(),
	}
	s.Server = httptest.NewUnstartedServer(s.Handler())

	return s, nil
}

func generateRepoName() string {
	size := 7
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
