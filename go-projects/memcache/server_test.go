package memcache

import (
	"testing"
)

func TestInvalidOperator(t *testing.T) {
	resp, err := makeTcpRequest([]byte("INVALID\r\n"))
	if resp != InvalidOperation.Error() {
		t.Fatalf(`makeTcpRequest("INVALID\r\n") = %q, %v, want %q, %v`, resp, err, InvalidOperation.Error(), nil)
	}
}
