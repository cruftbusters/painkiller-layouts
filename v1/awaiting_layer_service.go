package v1

import (
	"errors"

	"github.com/cruftbusters/painkiller-layouts/types"
)

type AwaitingLayerService interface {
	Enqueue(types.Layout) error
	Dequeue(priority int) types.Layout
}

var ErrQueueFull error = errors.New("error: queue full")

func NewAwaitingLayerService(queueSize int) AwaitingLayerService {
	return (&DefaultAwaitingLayerService{
		channel:       make(chan types.Layout, queueSize-1),
		addSubscriber: make(chan *struct{}),
		subscribers:   []Subscriber{},
	}).Launch()
}

type DefaultAwaitingLayerService struct {
	channel       chan types.Layout
	addSubscriber chan *struct{}
	subscribers   []Subscriber
}

type Subscriber struct {
	priority int
	channel  chan types.Layout
}

func (s *DefaultAwaitingLayerService) Launch() AwaitingLayerService {
	go func() {
		for {
			layout := <-s.channel
			if len(s.subscribers) == 0 {
				<-s.addSubscriber
			}
			subscriber, subscribers := s.subscribers[0], s.subscribers[1:]
			subscriber.channel <- layout
			s.subscribers = subscribers
		}
	}()
	return s
}

func (s *DefaultAwaitingLayerService) Enqueue(layout types.Layout) error {
	select {
	case s.channel <- layout:
		return nil
	default:
		return ErrQueueFull
	}
}

func (s *DefaultAwaitingLayerService) Dequeue(priority int) types.Layout {
	subscriber := Subscriber{priority, make(chan types.Layout)}
	s.subscribers = enqueueSubscriber(s.subscribers, subscriber)
	select {
	case s.addSubscriber <- nil:
	default:
	}
	return <-subscriber.channel
}

func enqueueSubscriber(subscribers []Subscriber, subscriber Subscriber) []Subscriber {
	for i := 0; i < len(subscribers); i++ {
		if subscriber.priority < subscribers[i].priority {
			subscribers = append(subscribers[:i+1], subscribers[i:]...)
			subscribers[i] = subscriber
			return subscribers
		}
	}
	return append(subscribers, subscriber)
}
