package snowflake

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ID int64
type SnowFlake struct {
	epoch        int64
	nodeBits     uint8
	sequenceBits uint8
	maxNodes     int64
	sequenceMask int64

	node      int64
	state     atomic.Int64 // timestamp + sequence
	startTime time.Time
	baseTime  int64
}
type Option func(*SnowFlake)

func WithEpoch(epoch int64) func(*SnowFlake) {
	return func(s *SnowFlake) {
		s.epoch = epoch
	}
}

func WithNodeBits(nodeBits uint8) func(*SnowFlake) {
	return func(s *SnowFlake) {
		s.nodeBits = nodeBits
	}
}

func WithSequenceBits(sequenceBits uint8) func(*SnowFlake) {
	return func(s *SnowFlake) {
		s.sequenceBits = sequenceBits
	}
}

// We get the current time from the monotonic clock
func (s *SnowFlake) now() int64 {
	return s.baseTime + time.Since(s.startTime).Milliseconds()
}

func NewSnowFlake(node int64, opt ...Option) *SnowFlake {

	s := SnowFlake{
		node:         node,
		nodeBits:     10,
		sequenceBits: 12,
	}

	for _, options := range opt {
		options(&s)
	}

	if s.nodeBits > 22 || s.sequenceBits > 41 || s.nodeBits+s.sequenceBits > 63 {
		panic("bits overflow int64")
	}

	var (
		epoch        int64 = 1288834974657
		maxNodes     int64 = (1 << s.nodeBits) - 1
		sequenceMask int64 = (1 << s.sequenceBits) - 1
	)

	if node < 0 || node > maxNodes {
		panic("snowflake id generator: Invalid node id")
	}

	if s.epoch == 0 {
		s.epoch = epoch
	}
	s.maxNodes = maxNodes
	s.sequenceMask = sequenceMask

	now := time.Now()
	base := now.UnixMilli() - s.epoch

	if base < 0 {
		panic("snowflake id generator: Invalid clock time")
	}

	s.state.Store(s.baseTime << int64(s.sequenceBits))
	s.startTime = now
	s.baseTime = base

	return &s
}

func (s *SnowFlake) Generate() ID {

	for {
		prev := s.state.Load()

		lastTime := prev >> s.sequenceBits
		seq := prev & s.sequenceMask

		now := s.now()
		if now < lastTime {
			now = lastTime
		}

		if now == lastTime {
			seq = (seq + 1) & s.sequenceMask
			if seq == 0 {
				now++
			}
		} else {
			seq = 0
		}

		next := now<<int64(s.sequenceBits) | seq

		if s.state.CompareAndSwap(prev, next) {
			var res int64 = now<<(int64(s.nodeBits)+int64(s.sequenceBits)) | s.node<<int64(s.sequenceBits) | seq
			return ID(res)
		}
	}
}

func (s *SnowFlake) GenerateN(n int) ([]ID, error) {

	if n < 0 || n > (1<<s.sequenceBits) {
		return nil, fmt.Errorf("Invalid n value, might cause collusion in IDs")
	}

	ids := make([]ID, n)
	i := 0

	for i < n {
		prev := s.state.Load()

		lastTime := prev >> s.sequenceBits
		seq := prev & s.sequenceMask

		now := s.now()
		if now < lastTime {
			now = lastTime
		}

		if now == lastTime {
			seq = (seq + 1) & s.sequenceMask
			if seq == 0 {
				now++
			}
		} else {
			seq = 0
		}

		next := now<<int64(s.sequenceBits) | seq

		if s.state.CompareAndSwap(prev, next) {
			var res int64 = now<<(int64(s.nodeBits)+int64(s.sequenceBits)) | s.node<<int64(s.sequenceBits) | seq
			ids[i] = ID(res)
			i++
		}
	}

	return ids, nil
}

// Node extracts the node ID from the Snowflake ID.
func (s *SnowFlake) Node(id ID) int64 {
	return (int64(id) >> int64(s.sequenceBits)) & s.maxNodes
}

// Timestamp extracts the relative timestamp (ms since epoch) from the Snowflake ID.
func (s *SnowFlake) Timestamp(id ID) int64 {
	return int64(id) >> (int64(s.sequenceBits) + int64(s.nodeBits))
}

// Sequence extracts the sequence number from the Snowflake ID.
func (s *SnowFlake) Sequence(id ID) int64 {
	return int64(id) & s.sequenceMask
}
