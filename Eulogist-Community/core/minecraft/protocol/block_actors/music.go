package block_actors

import (
	"Eulogist/core/minecraft/protocol"
	general "Eulogist/core/minecraft/protocol/block_actors/general_actors"
)

// 音符盒
type Music struct {
	general.BlockActor `mapstructure:",squash"`
	Note               byte `mapstructure:"note"` // TAG_Byte(1) = 0
}

// ID ...
func (*Music) ID() string {
	return IDMusic
}

func (n *Music) Marshal(io protocol.IO) {
	protocol.Single(io, &n.BlockActor)
	protocol.NBTInt(&n.Note, io.Varuint32)
}
