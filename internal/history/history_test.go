package history

import "testing"

func TestStoreKeepsNewestEntriesWithinLimit(t *testing.T) {
	store := New(2)
	store.Add("first")
	store.Add("second")
	store.Add("third")

	entries := store.Entries()

	if len(entries) != 2 {
		t.Fatalf("len(Entries()) = %d, want 2", len(entries))
	}
	if entries[0] != "second" || entries[1] != "third" {
		t.Errorf("Entries() = %q, want [second third]", entries)
	}
}

func TestStoreReturnsCopy(t *testing.T) {
	store := New(2)
	store.Add("first")

	entries := store.Entries()
	entries[0] = "changed"

	if store.Entries()[0] != "first" {
		t.Error("Entries() exposed internal state")
	}
}
