package containerd

import (
	"context"
	"errors"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/events"
	"github.com/containerd/containerd/namespaces"

	"github.com/weaveworks/reignite/core/ports"
	"github.com/weaveworks/reignite/pkg/defaults"
)

func NewEventService(cfg *Config) (ports.EventService, error) {
	client, err := containerd.New(cfg.SocketPath)
	if err != nil {
		return nil, fmt.Errorf("creating containerd client: %w", err)
	}

	return NewEventServiceWithClient(client), nil
}

func NewEventServiceWithClient(client *containerd.Client) ports.EventService {
	return &eventService{
		client: client,
	}
}

type eventService struct {
	client *containerd.Client
}

// Publish will publish an event to a specific topic.
func (es *eventService) Publish(ctx context.Context, topic string, eventToPublish interface{}) error {
	namespaceCtx := namespaces.WithNamespace(ctx, defaults.ContainerdNamespace)
	ctrEventSrv := es.client.EventService()
	if err := ctrEventSrv.Publish(namespaceCtx, topic, eventToPublish); err != nil {
		return fmt.Errorf("publishing event: %w", err)
	}

	return nil
}

// SubscribeTopic will subscribe to events on a named topic.
func (es *eventService) SubscribeTopic(ctx context.Context, topic string) (ch <-chan *ports.EventEnvelope, errs <-chan error) {
	topicFilter := fmt.Sprintf("topic==\"%s\"", topic)

	return es.subscribe(ctx, topicFilter)
}

// Subscribe will subscribe to events on all topics.
func (es *eventService) Subscribe(ctx context.Context) (ch <-chan *ports.EventEnvelope, errs <-chan error) {
	return es.subscribe(ctx)
}

func (es *eventService) subscribe(ctx context.Context, filters ...string) (ch <-chan *ports.EventEnvelope, errs <-chan error) {
	var (
		evtCh    = make(chan *ports.EventEnvelope)
		evtErrCh = make(chan error, 1)
	)
	errs = evtErrCh
	ch = evtCh

	namespaceCtx := namespaces.WithNamespace(ctx, defaults.ContainerdNamespace)

	var ctrEvents <-chan *events.Envelope
	var ctrErrs <-chan error
	if len(filters) == 0 {
		ctrEvents, ctrErrs = es.client.Subscribe(namespaceCtx)
	} else {
		ctrEvents, ctrErrs = es.client.Subscribe(namespaceCtx, filters...)
	}

	go func() {
		defer close(evtCh)

		for {
			select {
			case <-ctx.Done():
				if cerr := ctx.Err(); cerr != nil && !errors.Is(cerr, context.Canceled) {
					evtErrCh <- cerr
				}

				return
			case ctrEvt := <-ctrEvents:
				converted, err := convertCtrEventEnvelope(ctrEvt)
				if err != nil {
					evtErrCh <- fmt.Errorf("converting containerd event envelope: %w", err)
				}
				evtCh <- converted
			case ctrErr := <-ctrErrs:
				evtErrCh <- ctrErr
			}
		}
	}()

	return ch, errs
}