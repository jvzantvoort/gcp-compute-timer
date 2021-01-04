// Test code. Nothing useful trust me
package gcpcomputetimer

import (
	"testing"
)

func TestSecondsToHuman(t *testing.T) {
	results := make(map[int64]string)
	results[0] = "                 "
	results[1] = "         00:00:01"
	results[10] = "         00:00:10"
	results[100] = "         00:01:40"
	results[1000] = "         00:16:40"
	results[10000] = "         02:46:40"
	results[100000] = "  1 days 03:46:40"
	results[1000000] = " 11 days 13:46:40"

	fail := false
	for k, v := range results {
		r := SecondsToHuman(k)
		if r != v {
			t.Errorf("Expected \"%s\" got \"%s\"", v, r)
			fail = true
		}
	}
	if fail {
		t.Fail()
	}
}
