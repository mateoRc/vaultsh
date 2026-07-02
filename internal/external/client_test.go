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
		case "/summary":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(
				`{"requests":0,"errors":0,"avg_ms":0,"services":{},"commands":{}}`,
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
	if _, err := client.Summary(); err != nil {
		t.Fatal(err)
	}
}
