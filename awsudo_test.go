package awsudo

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilterPrefixes(t *testing.T) {
	testStrings := []string{"A", "AA", "AB", "BA", "BB", "B"}
	for tn, tc := range map[string]struct {
		prefixes []string
		want     []string
	}{
		"empty": {[]string{}, testStrings},
		"A":     {[]string{"A"}, []string{"BA", "BB", "B"}},
		"B":     {[]string{"B"}, []string{"A", "AA", "AB"}},
		"C":     {[]string{"C"}, testStrings},
	} {
		t.Run(tn, func(t *testing.T) {
			got := filterPrefixes(testStrings, tc.prefixes)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("mismatch: (-got,+want):\n%v", diff)
			}
		})
	}
}
