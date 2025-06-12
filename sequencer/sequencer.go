package sequencer

import (
	"abchain_scan/config"
	"abchain_scan/log"
	"go.uber.org/zap"
	"sync"
)

type Sequenceable interface {
	GetSequence() uint64
}

type Committable interface {
	Commit(value Sequenceable)
}

type Sequencer interface {
	Init(height uint64)
	CommitWithSequence(value Sequenceable, output Committable)
}

type sequencer struct {
	active   bool
	mu       sync.Mutex
	cond     *sync.Cond
	sequence uint64
}

func NewSequencer() Sequencer {
	s := &sequencer{
		active: config.G.EnableSequencer,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *sequencer) Init(sequence uint64) {
	log.Logger.Info("init sequencer", zap.Uint64("sequence", sequence))
	if s.sequence == 0 {
		s.sequence = sequence - 1
	} else {
		log.Logger.Fatal("sequencer init err", zap.Uint64("sequence", sequence), zap.Uint64("old sequence", s.sequence))
	}
}

func (s *sequencer) CommitWithSequence(value Sequenceable, output Committable) {
	if !s.active {
		output.Commit(value)
		return
	}

	sequence := value.GetSequence()

	s.mu.Lock()
	for s.sequence+1 != sequence {
		s.cond.Wait()
	}

	output.Commit(value)
	s.sequence = sequence

	s.cond.Broadcast()
	s.mu.Unlock()
}
