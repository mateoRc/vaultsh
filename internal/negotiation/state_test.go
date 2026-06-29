package negotiation

import "testing"

func TestStateAcceptsPendingOfferOnce(t *testing.T) {
	state := &State{}
	state.Propose(150000)

	offer, found := state.Accept()
	if !found {
		t.Fatal("Accept() did not find pending offer")
	}
	if offer != 150000 {
		t.Errorf("Accept() offer = %f, want 150000", offer)
	}

	if _, found := state.Accept(); found {
		t.Error("Accept() reused an accepted offer")
	}
}
