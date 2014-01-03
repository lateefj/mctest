package mctest

import (
  "fmt"
  "net/http"
  "testing"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "HomeHandler")
}

func ByteHomeHandler(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("HomeHandler"))
}

func TestHome(t *testing.T) {
  req, _ := http.NewRequest("GET", "/path/to/handler", nil)
  resp := NewMockTestResponse(t)
  HomeHandler(resp, req)
  resp.AssertBody("HomeHandler")
  resp.AssertCode(200)
}

func TestByteHome(t *testing.T) {
  req, _ := http.NewRequest("GET", "/path/to/handler", nil)
  resp := NewMockTestResponse(t)
  ByteHomeHandler(resp, req)
  resp.AssertCode(200)
  if string(resp.Bytes()) != "HomeHandler" {
    t.Fatalf("Expected bytes to equal HomeHandler but failed actual is '%s'", string(resp.Bytes()))
  }
  resp.AssertBody("HomeHandler")
}
