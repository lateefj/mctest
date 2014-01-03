package mctest

import (
  "bytes"
  "net/http"
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
func (mtr *MockTestResponse) AssertBody(b string) {
  if mtr.String() != b {
    mtr.T.Fatalf("Response body is %s asserted that it is %s", mtr.String(), b)
  }
}

// AssertCode ... Helper method to validate the right status code
func (mtr *MockTestResponse) AssertCode(c int) {
  if mtr.StatusCode != c {
    mtr.T.Fatalf("Response StatusCode is %d asserted that it is %d", mtr.StatusCode, c)
  }
}

// AssertHeaders ... Helper method to validate the proper headers retuned
func (mtr *MockTestResponse) AssertHeaders(t *testing.T, expectedHeaders map[string]string) {
  if len(expectedHeaders) != len(mtr.Head) {
    mtr.T.Fatalf("expected %v headers; got %v. Returned headers: %v", len(expectedHeaders), len(mtr.Head), mtr.Head)
  }
  for key := range expectedHeaders {
    if val := mtr.Head.Get(key); val != expectedHeaders[key] {
      mtr.T.Fatalf("expected header %v to be '%v'; got '%v'", key, expectedHeaders[key], val)
    }
  }
}

// NewMockTestResponse ... Create an instance of MockTestResponse
func NewMockTestResponse(t *testing.T) *MockTestResponse {
  return &MockTestResponse{MockResponse: MockResponse{Buffer: *bytes.NewBuffer(make([]byte, 0)), Head: http.Header{}, StatusCode: -1}, T: t}
}
