package pkg

import (
	"sync"
	"time"
)

type Generator interface {
	ID() int64
}

var (
	// first bit is 0
	timestampLength  int64 = 41
	instanceIDLength int64 = 10
	seedLength       int64 = 12
	timestampOffset  int64 = instanceIDLength + seedLength
	instanceIDOffset int64 = seedLength
	seedOffset       int64 = 0
	timestampMask    int64 = ((1 << timestampLength) - 1) << timestampOffset
	instanceIDMask   int64 = ((1 << instanceIDLength) - 1) << instanceIDOffset
	seedMask         int64 = ((1 << seedLength) - 1) << seedOffset

	generatorInit sync.Once
	generator     Generator
)

type Snowflake int64

func (s Snowflake) Timestamp() int64 {
	return (int64(s) & timestampMask) >> timestampOffset
}

func (s Snowflake) InstanceID() uint16 {
	return uint16((int64(s) & instanceIDMask) >> instanceIDOffset)
}

func (s Snowflake) Seed() uint16 {
	return uint16((int64(s) & seedMask) >> seedOffset)
}

func (s Snowflake) Encode() string {
	return EncodeBase62(int64(s))
}

func ParseBase62(s string) (Snowflake, error) {
	id, err := DecodeBase62(s)
	return Snowflake(id), err
}

type sfGenerator struct {
	mux           sync.Mutex
	epoch         int64
	instanceID    uint16
	seed          uint8
	lastResetTime int64
}

func NewSfGen(instanceID uint16) Generator {
	generatorInit.Do(func() {
		epoch, _ := time.Parse("2006-01-02", "2025-01-01")
		generator = &sfGenerator{
			epoch:         epoch.UTC().UnixMilli(),
			instanceID:    instanceID,
			seed:          0,
			lastResetTime: 0,
		}
	})
	return generator
}

func (g *sfGenerator) ID() int64 {
	g.mux.Lock()
	defer g.mux.Unlock()
	unixTimestamp := time.Now().UTC().UnixMilli() / 10
	timestamp := unixTimestamp - (g.epoch / 10)
	if unixTimestamp-g.lastResetTime >= 1 {
		g.seed = 0
		g.lastResetTime = unixTimestamp
	}
	g.seed++
	id := int64(0)
	id |= ((timestamp << timestampOffset) & timestampMask)
	id |= ((int64(g.instanceID) << instanceIDOffset) & instanceIDMask)
	id |= ((int64(g.seed) << seedOffset) & seedMask)
	return id
}
