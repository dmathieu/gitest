package gitest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitServer(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()
	assert.NotNil(t, server)
}

func TestCallHomePage(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	res, err := http.Get(server.URL)
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func TestCallRefs(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	url := fmt.Sprintf("%s/%s.git/info/refs?service=git-upload-pack", server.URL, server.ValidRepo)
	res, err := http.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Regexp(t, regexp.MustCompile("service=git-upload-pack"), string(content))
}

func TestCallRefsInvalidService(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	url := fmt.Sprintf("%s/%s.git/info/refs?service=ls", server.URL, server.ValidRepo)
	res, err := http.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

func TestCallRefsUnknownRepo(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	res, err := http.Get(fmt.Sprintf("%s/%s.git/info/refs", server.URL, "unknown_repo"))
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func TestCallRefsNotAllowedRepo(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	res, err := http.Get(fmt.Sprintf("%s/%s.git/info/refs", server.URL, server.NotAllowedRepo))
	assert.Nil(t, err)
	assert.Equal(t, 401, res.StatusCode)
}

func TestCallService(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	res, err := http.Post(fmt.Sprintf("%s/%s.git/git-upload-pack", server.URL, server.ValidRepo), "text/plain", &bytes.Buffer{})
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestCallServiceUnknownRepo(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	res, err := http.Post(fmt.Sprintf("%s/%s.git/git-upload-pack", server.URL, "unknown_repo"), "text/plain", &bytes.Buffer{})
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func TestCallServiceNotAllowedRepo(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	res, err := http.Post(fmt.Sprintf("%s/%s.git/git-upload-pack", server.URL, server.NotAllowedRepo), "text/plain", &bytes.Buffer{})
	assert.Nil(t, err)
	assert.Equal(t, 401, res.StatusCode)
}

func TestGitClone(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	tempDir, err := ioutil.TempDir("", "git_repository")
	assert.Nil(t, err)

	var out = new(bytes.Buffer)
	c := exec.Command("git", "clone", fmt.Sprintf("%s/%s.git", server.URL, server.ValidRepo), tempDir)
	c.Stdout = out
	c.Stderr = out
	err = c.Run()
	assert.Nil(t, err)
	assert.Regexp(t, regexp.MustCompile("Cloning into"), out)
}

func TestGitPush(t *testing.T) {
	server, err := NewServer("basic")
	assert.Nil(t, err)
	defer server.Close()

	tempDir, err := ioutil.TempDir("", "git_repository")
	assert.Nil(t, err)

	err = exec.Command("git", "clone", fmt.Sprintf("%s/%s.git", server.URL, server.ValidRepo), tempDir).Run()
	assert.Nil(t, err)

	cmd := exec.Command("git", "commit", "-m", "test", "--allow-empty")
	cmd.Dir = tempDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	assert.Nil(t, err)

	cmd = exec.Command("git", "push", "--no-verify")
	cmd.Dir = tempDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	assert.Nil(t, err)
}
