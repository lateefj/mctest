package mctest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "HomeHandler")
}

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

// TODO: Write some more complex bytes
func ByteHomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("HomeHandler"))
}
func TestByteHome(t *testing.T) {
	req, _ := http.NewRequest("GET", "/path/to/handler", nil)
	resp := NewMockTestResponse(t)
	ByteHomeHandler(resp, req)
	b := "HomeHandler"
	if string(resp.Bytes()) != "HomeHandler" {
		t.Fatalf("Expected bytes to equal HomeHandler but failed actual is '%s'", string(resp.Bytes()))
	}
	if !resp.AssertCode(200) {
		t.Fatalf("Response StatusCode is %d asserted that it is %d", resp.StatusCode, 200)
	}
	if !resp.AssertBody(b) {
		t.Fatalf("Response body is %s asserted that it is %s", resp.String(), b)
	}
}

type Payload struct {
	X string `json:"x"`
	Y string `json:"y"`
}

func testJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	p := &Payload{X: "foo", Y: "bar"}
	j, err := json.Marshal(p)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error parsing json %s", err)))
	}
	w.Write(j)
}
func TestJsonAssert(t *testing.T) {
	req, _ := http.NewRequest("GET", "/path/to/handler", nil)
	nt := &testing.T{}
	resp := NewMockTestResponse(nt)
	testJsonHandler(resp, req)
	// Payload mathces so keep moving..
	if !resp.AssertJson(&Payload{}, &Payload{X: "foo", Y: "bar"}) {
		t.Fatalf("Failed to validate payload")
	}
	// Should fail
	if resp.AssertJson(&Payload{}, &Payload{X: "boo", Y: "bar"}) {
		t.Fatalf("Payload does not match!")
	}

}
