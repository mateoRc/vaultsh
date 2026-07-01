package telemetry

import (
	"context"
	"testing"
	"time"
)

type senderStub struct {
	received chan Event
	block    chan struct{}
}

func (s *senderStub) Record(
	service string,
	event string,
	name string,
	durationMS int64,
	exitCode int,
) error {
	if s.block != nil {
		<-s.block
	}
	s.received <- Event{
		Service:    service,
		Operation:  event,
		Name:       name,
		DurationMS: durationMS,
		ExitCode:   exitCode,
	}
	return nil
}

func TestDispatcherDeliversInBackground(t *testing.T) {
	sender := &senderStub{received: make(chan Event, 1)}
	dispatcher := NewDispatcher(sender, 1, nil)
	defer dispatcher.Close(context.Background())

	_ = dispatcher.Record("vault", "command.executed", "grep", 8, 0)

	select {
	case event := <-sender.received:
		if event.Name != "grep" || event.DurationMS != 8 {
			t.Errorf("event = %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("event was not delivered")
	}
}

func TestDispatcherDropsWhenQueueIsFull(t *testing.T) {
	unblock := make(chan struct{})
	sender := &senderStub{
		received: make(chan Event, 2),
		block:    unblock,
	}
	dispatcher := NewDispatcher(sender, 1, nil)

	_ = dispatcher.Record("vault", "command.executed", "first", 1, 0)
	time.Sleep(10 * time.Millisecond)
	_ = dispatcher.Record("vault", "command.executed", "second", 1, 0)
	_ = dispatcher.Record("vault", "command.executed", "dropped", 1, 0)

	if dispatcher.Dropped() != 1 {
		t.Errorf("Dropped() = %d, want 1", dispatcher.Dropped())
	}

	close(unblock)
	dispatcher.Close(context.Background())
}
