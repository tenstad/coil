package coild

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testGetStatus(t *testing.T) {
	t.Parallel()
	mockDB := newMock()
	server := NewServer(mockDB)
	server.containerIPs = map[string][]net.IP{
		"container-1": {net.ParseIP("10.0.0.1")},
	}

	_, subnet1, _ := net.ParseCIDR("10.0.0.0/27")
	server.addressBlocks = []*net.IPNet{subnet1}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/status", nil)
	server.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("http status should be 200, actual:", resp.StatusCode)
	}
	st := status{}
	err := json.NewDecoder(resp.Body).Decode(&st)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(st.AddressBlocks, []string{"10.0.0.0/27"}) {
		t.Error(`expected: []string{"10.0.0.0/27"}, actual:`, st.AddressBlocks)
	}
	if !cmp.Equal(st.Containers, map[string][]string{
		"container-1": {"10.0.0.1"},
	}) {
		t.Error(`expected: "container-1": {"10.0.0.1"}, actual:`, st.Containers)
	}
	if st.Status != http.StatusOK {
		t.Error("expected: 200, actual:", st.Status)
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/status", nil)
	server.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Error(`should not be 200`)
	}
}

func TestServeHTTP(t *testing.T) {
	t.Run("status", testGetStatus)
}
