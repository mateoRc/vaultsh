package external

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchSendsAtlasBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer atlas-secret" {
			t.Errorf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "atlas-secret", "")
	if _, err := client.Search("kafka"); err != nil {
		t.Fatal(err)
	}
}

func TestForgeRequestsSendBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer forge-secret" {
			t.Errorf("Authorization = %q", got)
		}
		switch r.URL.Path {
		case "/events":
			w.WriteHeader(http.StatusNoContent)
		case "/visits":
			w.WriteHeader(http.StatusNoContent)
		case "/summary":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(
				`{"requests":0,"errors":0,"user_errors":0,"avg_ms":0,"services":{},"commands":{}}`,
			))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient("", server.URL, "", "forge-secret")
	if err := client.Record("vault", "command.executed", "ls", 1, 0); err != nil {
		t.Fatal(err)
	}
	if err := client.RecordVisit("12345678-1234-4123-8123-123456789abc"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Summary(); err != nil {
		t.Fatal(err)
	}
}

func TestDashboardSendsWindowAndInvalidFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dashboard" || r.URL.Query().Get("window") != "all" ||
			r.URL.Query().Get("include_invalid") != "true" ||
			r.Header.Get("Authorization") != "Bearer forge-secret" {
			t.Errorf("unexpected dashboard request: %s", r.URL.String())
		}
		_, _ = w.Write([]byte("ALL TIME"))
	}))
	defer server.Close()

	client := NewClient("", server.URL, "", "forge-secret")
	output, err := client.Dashboard(0, true)
	if err != nil || output != "ALL TIME" {
		t.Fatalf("dashboard = %q, %v", output, err)
	}
}
