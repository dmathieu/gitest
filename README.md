# Gitest

[![Build Status](https://travis-ci.org/dmathieu/gitest.svg?branch=master)](https://travis-ci.org/dmathieu/gitest)

A mock GIT server for running `git clone` in your app's tests without relying
on the network

## Usage

```
server, err := gitest.NewServer("basic")
if err != nil {
  log.Fatalf(err)
}
defer server.Close()

tempDir, err := ioutil.TempDir("", "git_repository")
if err != nil {
  log.Fatalf(err)
}

c := exec.Command("git", "clone", fmt.Sprintf("%s/%s.git", server.URL, server.ValidRepo), tempDir)
err = c.Run()
if err != nil {
  log.Fatalf(err)
}
```
