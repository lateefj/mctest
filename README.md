# McTest

![Northern Mockingbird](http://www.biokids.umich.edu/collections/contributors/grzimek_birds/Mimidae/Mimus_polyglottos/medium.jpg)

[Image Source](http://www.biokids.umich.edu/critters/Mimus_polyglottos/pictures/resources/contributors/grzimek_birds/Mimidae/Mimus_polyglottos/)


[![Build Status](https://travis-ci.org/lateefj/mctest.svg?branch=master)](https://travis-ci.org/lateefj/mctest)

McTest is a Go (golang) web test library. Initally started as a gist and is slowly getting used more for testing webservices written in Go.


## Goals

 * Simple quick way to test Middleware (wrappers around http handler)
 * Quickly be able to test web service (JSON)
 * Highlevel function integration tests 

## Basic API

### AssertCode(http.StatusOK) bool

This simple failes the test if the code does not match the one passed in

### AssertBody(string)  bool

If the string that is passed in doesn't match then fail the test.

#### AssertJson(instance interface{}, expectedMatchStruct interface{}) bool

By passing a struct this will validate that the response json has the same values

#### Bytes() []byte

Comparing bytes can be achieved by calling resp.Bytes().


## Use Cases

### Basic

```go
// Application code example
func HomeHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "HomeHandler")
}

// Testing application code
func TestHome(t *testing.T) {
  req, _ := http.NewRequest("GET", "/path/to/handler", nil)
  resp := NewMockTestResponse(t)
  HomeHandler(resp, req)
  b := "HomeHandler"
  if !resp.AssertCode(http.StatusOK) {
    t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, http.StatusOK)
  }
  if !resp.AssertBody(b) {
    t.Fatalf("Response body is %s asserted that it is %s", resp.String(), b)
  }
}
```

### Json

```go

// JSON Stuct
type payload struct {
  X string `json:"x"`
  Y string `json:"y"`
}
// JSON API
func HomeHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "{\"x\":\"bar\", \"y\":\"foo\"}")
}
// Test JSON API
func TestHome(t *testing.T) {
  req, _ := http.NewRequest("GET", "/path/to/handler", nil)
  resp := NewMockTestResponse(t)
  HomeHandler(resp, req)
  b := "HomeHandler"
  if !resp.AssertCode(http.StatusOK) {
    t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, http.StatusOK)
  }
  p := payload{X: "bar", Y: "foo"}
  inst := &payload{}
  if !resp.AssertJson(inst, p) {
    t.Fatalf("Response body is %s asserted that it is %v", resp.String(), p)
  }
}
```

### Middleware

```go
// App authentication wrapper
func testAuthWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a, err := r.Cookie(authCookieName)
		if err != nil || a == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(failedAuth))
		} else {
			fn(w, r)
		}
	}
}
// Handler code
func testAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("TestAuthHandler"))
}
// Actual test code
func TestAuthMiddleware(t *testing.T) {
	req, _ := http.NewRequest("GET", "/path/to/handler", nil)
  // Create mock
	resp := NewMockTestResponse(t)
	testAuthWrapper(testAuthHandler)(resp, req)
	if !resp.AssertCode(http.StatusUnauthorized) {
		t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, http.StatusUnauthorized)
	}
	if !resp.AssertBody(failedAuth) {
		t.Fatalf("Response body is %s asserted that it is %s", resp.String(), failedAuth)
	}
	b := "TestAuthHandler"
  // Add auth cookie
	req.AddCookie(authCookie)
	resp = NewMockTestResponse(t)
	// Run the code again
	testAuthWrapper(testAuthHandler)(resp, req)
	if !resp.AssertCode(http.StatusOK) {
		t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, http.StatusOK)
	}
	if !resp.AssertBody(b) {
		t.Fatalf("Response body is %s asserted that it is %s", resp.String(), b)
	}
}
```

Check the http_test.go for more examples on how to use the api.

