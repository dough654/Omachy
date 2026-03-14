package preflight

import "testing"

func TestAllPassedAllGreen(t *testing.T) {
	checks := []Check{
		{Name: "A", Passed: true},
		{Name: "B", Passed: true},
		{Name: "C", Passed: true},
	}
	if !AllPassed(checks) {
		t.Error("all checks passed, AllPassed should return true")
	}
}

func TestAllPassedWithWarning(t *testing.T) {
	checks := []Check{
		{Name: "A", Passed: true},
		{Name: "B", Passed: false, Warning: true},
		{Name: "C", Passed: true},
	}
	if !AllPassed(checks) {
		t.Error("warning-only failure should not cause AllPassed to return false")
	}
}

func TestAllPassedWithFailure(t *testing.T) {
	checks := []Check{
		{Name: "A", Passed: true},
		{Name: "B", Passed: false, Warning: false},
		{Name: "C", Passed: true},
	}
	if AllPassed(checks) {
		t.Error("non-warning failure should cause AllPassed to return false")
	}
}

func TestMacOSName(t *testing.T) {
	tests := []struct {
		major int
		want  string
	}{
		{13, "Ventura"},
		{14, "Sonoma"},
		{15, "Sequoia"},
		{16, "Tahoe"},
		{99, "macOS"},
	}
	for _, tt := range tests {
		got := macOSName(tt.major)
		if got != tt.want {
			t.Errorf("macOSName(%d) = %q, want %q", tt.major, got, tt.want)
		}
	}
}
