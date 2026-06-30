package command

import "testing"

func TestWhoami(t *testing.T) {
	result := (Whoami{}).Execute(nil, Input{})
	want := "Mateo Mahmutović\n" +
		"Senior Backend Engineer\n" +
		"Currently building distributed backend systems."

	if result.Output != want {
		t.Errorf("output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != ExitSuccess {
		t.Errorf("exit code = %d, want %d", result.ExitCode, ExitSuccess)
	}
}

func TestWhoamiIsHidden(t *testing.T) {
	if !IsHidden(Whoami{}) {
		t.Error("whoami should be hidden from the general help listing")
	}
}
