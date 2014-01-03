McTest
======

McTest is a Go (golang) web test library. Initally started as a gist and is slowly getting used more for testing webservices written in Go.



```go

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(200)
  fmt.Fprintf(w, "HomeHandler")
}

func TestHome(t *testing.T) {
  req, _ := http.NewRequest("GET", "/path/to/handler", nil)
  resp := NewMockTestResponse(t)
  HomeHandler(resp, req)
  b := "HomeHandler"
  if !resp.AssertCode(200) {
    t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, 200)
  }
  if !resp.AssertBody(b) {
    t.Fatalf("Response body is %s asserted that it is %s", resp.String(), b)
  }
}
```

The interesting functions are:

### AssertCode
This simple failes the test if the code does not match the one passed in

### AssertBody
If the string that is passed in doesn't match then fail the test.

#### Bytes()
Comparing bytes can be achieved by calling resp.Bytes().


