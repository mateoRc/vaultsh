package telemetry

import (
	"context"
	"log/slog"
	"sync/atomic"
)

type Event struct {
	Service    string
	Name       string
	Operation  string
	DurationMS int64
	ExitCode   int
}

type Sender interface {
	Record(service, event, name string, durationMS int64, exitCode int) error
}

type Dispatcher struct {
	events  chan Event
	sender  Sender
	logger  *slog.Logger
	stop    chan struct{}
	stopped chan struct{}
	dropped atomic.Uint64
}

func NewDispatcher(sender Sender, capacity int, logger *slog.Logger) *Dispatcher {
	if capacity <= 0 {
		capacity = 1000
	}
	dispatcher := &Dispatcher{
		events:  make(chan Event, capacity),
		sender:  sender,
		logger:  logger,
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
	go dispatcher.run()
	return dispatcher
}

func (d *Dispatcher) Record(
	service string,
	event string,
	name string,
	durationMS int64,
	exitCode int,
) error {
	item := Event{
		Service:    service,
		Operation:  event,
		Name:       name,
		DurationMS: durationMS,
		ExitCode:   exitCode,
	}

	select {
	case <-d.stop:
		d.drop()
	case d.events <- item:
	default:
		d.drop()
	}
	return nil
}

func (d *Dispatcher) Close(ctx context.Context) {
	select {
	case <-d.stop:
	default:
		close(d.stop)
	}

	select {
	case <-d.stopped:
	case <-ctx.Done():
	}
}

func (d *Dispatcher) Dropped() uint64 {
	return d.dropped.Load()
}

func (d *Dispatcher) run() {
	defer close(d.stopped)

	for {
		select {
		case item := <-d.events:
			d.send(item)
		case <-d.stop:
			for {
				select {
				case item := <-d.events:
					d.send(item)
				default:
					return
				}
			}
		}
	}
}

func (d *Dispatcher) send(item Event) {
	_ = d.sender.Record(
		item.Service,
		item.Operation,
		item.Name,
		item.DurationMS,
		item.ExitCode,
	)
}

func (d *Dispatcher) drop() {
	count := d.dropped.Add(1)
	if d.logger != nil {
		d.logger.Warn("telemetry event dropped", "dropped", count)
	}
}
