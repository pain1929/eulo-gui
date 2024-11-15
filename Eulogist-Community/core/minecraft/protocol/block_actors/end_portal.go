package block_actors

import general "Eulogist/core/minecraft/protocol/block_actors/general_actors"

// 末地折跃门
type EndPortal struct {
	general.BlockActor `mapstructure:",squash"`
}

// ID ...
func (*EndPortal) ID() string {
	return IDEndPortal
}
