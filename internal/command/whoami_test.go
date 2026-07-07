package command

import "testing"

func TestWhoami(t *testing.T) {
	result := (Whoami{}).Execute(nil, Input{})
	want := "Mateo Mahmutović\n" +
		"Senior Backend Engineer\n" +
		"Currently building distributed backend systems.\n\n" +
		"[Email](mailto:mahmutovic.mateo@gmail.com)\n" +
		"[GitHub](https://github.com/mateoRc)\n" +
		"[LinkedIn](https://www.linkedin.com/in/mateo-mahmutovi%C4%87-a9837232b/)"

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
