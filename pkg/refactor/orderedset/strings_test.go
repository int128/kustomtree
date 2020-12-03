package orderedset

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStrings_Append(t *testing.T) {
	t.Run("ZeroValue", func(t *testing.T) {
		var s Strings
		if s.Get() != nil {
			t.Errorf("Get() wants nil but was %+v", s.elements)
		}
	})
	t.Run("Dedupe", func(t *testing.T) {
		var s Strings
		s.Append("foo")
		if diff := cmp.Diff([]string{"foo"}, s.elements); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
		s.Append("foo")
		if diff := cmp.Diff([]string{"foo"}, s.elements); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
		s.Append("bar")
		if diff := cmp.Diff([]string{"foo", "bar"}, s.elements); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
		s.Append("foo")
		if diff := cmp.Diff([]string{"foo", "bar"}, s.elements); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
}
