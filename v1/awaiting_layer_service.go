package v1

import (
	"errors"

	"github.com/cruftbusters/painkiller-layouts/types"
)

type AwaitingLayerService interface {
	Enqueue(types.Layout) error
	Dequeue() types.Layout
}

var ErrQueueFull error = errors.New("queue full")

func NewAwaitingLayerService(queueSize int) AwaitingLayerService {
	return &DefaultAwaitingLayerService{
		channel: make(chan types.Layout, queueSize),
	}
}

type DefaultAwaitingLayerService struct {
	channel chan types.Layout
}

func (s *DefaultAwaitingLayerService) Enqueue(layout types.Layout) error {
	select {
	case s.channel <- layout:
		return nil
	default:
		return ErrQueueFull
	}
}

func (s *DefaultAwaitingLayerService) Dequeue() types.Layout {
	return <-s.channel
}
