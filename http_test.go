package mctest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

var (
	failedAuth      = "Failed to authenticate"
	authCookieName  = "auth_cookie"
	authCookieValue = "TEST_AUTH"
	authCookie      = &http.Cookie{Name: authCookieName, Value: authCookieValue}
)

func TestInit(t *testing.T) {
	req, _ := http.NewRequest("GET", "/path/to/handler", nil)
	resp := NewMockTestResponse(t)
	go func(w http.ResponseWriter, r *http.Request) {

	}(resp, req)
	if !resp.AssertCode(-1) {
		t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, -1)
	}
	if len(resp.Bytes()) != 0 {
		t.Fatalf("Expected response bytes to be 0 but they are %d", len(resp.Bytes()))
	}

}

func TestHome(t *testing.T) {
	req, _ := http.NewRequest("GET", "/path/to/handler", nil)
	resp := NewMockTestResponse(t)
	b := "HomeHandler"
	func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, b)
	}(resp, req)
	if !resp.AssertCode(http.StatusOK) {
		t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, http.StatusOK)
	}
	if !resp.AssertBody(b) {
		t.Fatalf("Response body is %s asserted that it is %s", resp.String(), b)
	}
}

func TestByteHome(t *testing.T) {
	req, _ := http.NewRequest("GET", "/path/to/handler", nil)
	resp := NewMockTestResponse(t)
	b := "HomeHandler"
	func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(b))
	}(resp, req)
	if string(resp.Bytes()) != "HomeHandler" {
		t.Fatalf("Expected bytes to equal HomeHandler but failed actual is '%s'", string(resp.Bytes()))
	}
	if !resp.AssertCode(http.StatusOK) {
		t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, http.StatusOK)
	}
	if !resp.AssertBody(b) {
		t.Fatalf("Response body is %s asserted that it is %s", resp.String(), b)
	}
}

type Payload struct {
	X string `json:"x"`
	Y string `json:"y"`
}

func TestJsonAssert(t *testing.T) {
	req, _ := http.NewRequest("GET", "/path/to/handler", nil)
	nt := &testing.T{}
	resp := NewMockTestResponse(nt)
	func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		p := &Payload{X: "foo", Y: "bar"}
		j, err := json.Marshal(p)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error parsing json %s", err)))
		}
		w.Write(j)
	}(resp, req)
	// Payload mathces so keep moving..
	if !resp.AssertJson(&Payload{}, &Payload{X: "foo", Y: "bar"}) {
		t.Fatalf("Failed to validate payload")
	}
	// Should fail
	if resp.AssertJson(&Payload{}, &Payload{X: "boo", Y: "bar"}) {
		t.Fatalf("Payload does not match!")
	}
}

// This validates that auth middleware does the right thing
func TestAuthMiddleware(t *testing.T) {
	req, _ := http.NewRequest("GET", "/path/to/handler", nil)
	resp := NewMockTestResponse(t)
	authWrapper := func(fn http.HandlerFunc) http.HandlerFunc {
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
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TestAuthHandler"))
	}
	authWrapper(handler)(resp, req)
	if !resp.AssertCode(http.StatusUnauthorized) {
		t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, http.StatusUnauthorized)

	}

	if !resp.AssertBody(failedAuth) {
		t.Fatalf("Response body is %s asserted that it is %s", resp.String(), failedAuth)
	}

	b := "TestAuthHandler"
	req.AddCookie(authCookie)
	resp = NewMockTestResponse(t)
	// Run the code again
	authWrapper(handler)(resp, req)
	if !resp.AssertCode(http.StatusOK) {
		t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, http.StatusOK)
	}
	if !resp.AssertBody(b) {
		t.Fatalf("Response body is %s asserted that it is %s", resp.String(), b)
	}

}
