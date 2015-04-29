package mctest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

// MockResponse [...] Create a mock http response implementaiton
type MockResponse struct {
	bytes.Buffer
	Head       http.Header
	StatusCode int
}

// Header ...  Returns the header
func (mr *MockResponse) Header() http.Header {
	return mr.Head
}

// WriteHeader ... Write a response code to the header
func (mr *MockResponse) WriteHeader(code int) {
	mr.StatusCode = code
}

// MockTestResponse [...] Simple wrapper to save typing
type MockTestResponse struct {
	MockResponse
	T *testing.T
}

// AssertBody ... Wrapper to provide an typing savings
func (mtr *MockTestResponse) AssertBody(b string) bool {
	if mtr.String() != b {
		mtr.T.Errorf("Response body is %s asserted that it is %s", mtr.String(), b)
		return false
	}
	return true
}

// AssertCode ... Helper method to validate the right status code
func (mtr *MockTestResponse) AssertCode(c int) bool {
	if mtr.StatusCode != c {
		mtr.T.Errorf("Response StatusCode is %d asserted that it is %d", mtr.StatusCode, c)
		return false
	}
	return true
}

// AssertHeaders ... Helper method to validate the proper headers retuned
func (mtr *MockTestResponse) AssertHeaders(expectedHeaders map[string]string) {
	if len(expectedHeaders) != len(mtr.Head) {
		mtr.T.Fatalf("expected %v headers; got %v. Returned headers: %v", len(expectedHeaders), len(mtr.Head), mtr.Head)
	}
	for key := range expectedHeaders {
		if val := mtr.Head.Get(key); val != expectedHeaders[key] {
			mtr.T.Fatalf("expected header %v to be '%v'; got '%v'", key, expectedHeaders[key], val)
		}
	}
}

func (mtr *MockTestResponse) AssertJson(inst, data interface{}) bool {
	err := json.Unmarshal(mtr.Bytes(), inst)
	if err != nil {
		mtr.T.Errorf("Failed to unmarshal with error: %s", err)
	}
	if !reflect.DeepEqual(inst, data) {
		// Try to extract json text for displaying error
		txt, err := json.Marshal(data)
		if err != nil {
			mtr.T.Errorf("Type %s did not match response %v: expected %v", reflect.TypeOf(data).Kind().String(), mtr.String(), data)
			return false
		}
		mtr.T.Errorf("Type %s did not match response %s: expected %s", reflect.TypeOf(data).Kind().String(), mtr.String(), txt)
		return false
	}
	return true

}

// NewMockTestResponse ... Create an instance of MockTestResponse
func NewMockTestResponse(t *testing.T) *MockTestResponse {
	return &MockTestResponse{MockResponse: MockResponse{Buffer: *bytes.NewBuffer(make([]byte, 0)), Head: http.Header{}, StatusCode: -1}, T: t}
}
