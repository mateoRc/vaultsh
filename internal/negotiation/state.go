package negotiation

import "sync"

type State struct {
	mu           sync.Mutex
	pending      bool
	counterOffer float64
}

func (s *State) Propose(counterOffer float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pending = true
	s.counterOffer = counterOffer
}

func (s *State) Accept() (float64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.pending {
		return 0, false
	}

	counterOffer := s.counterOffer
	s.pending = false
	s.counterOffer = 0
	return counterOffer, true
}
