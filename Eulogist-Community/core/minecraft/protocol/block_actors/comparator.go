package block_actors

import (
	"Eulogist/core/minecraft/protocol"
	general "Eulogist/core/minecraft/protocol/block_actors/general_actors"
)

// 比较器
type Comparator struct {
	general.BlockActor `mapstructure:",squash"`
	OutputSignal       int32 `mapstructure:"OutputSignal"` // TAG_Int(4) = 0
}

// ID ...
func (*Comparator) ID() string {
	return IDComparator
}

func (c *Comparator) Marshal(io protocol.IO) {
	protocol.Single(io, &c.BlockActor)
	io.Varint32(&c.OutputSignal)
}
