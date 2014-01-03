mctest
======

McTest is a Go (golang) web test library. Initally started as a gist and is slowly getting used more for testing webservices written in Go.



```go

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "HomeHandler")
}

func TestHome(t *testing.T) {
  req, _ := http.NewRequest("GET", "/path/to/handler", nil)
  resp := NewMockTestResponse(t)
  HomeHandler(resp, req)
  resp.AssertBody("HomeHandler")
  resp.AssertCode(200)
}
```


