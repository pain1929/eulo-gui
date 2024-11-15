package block_actors

import general "Eulogist/core/minecraft/protocol/block_actors/general_actors"

// 荧光物品展示框
type GlowItemFrame struct {
	general.ItemFrameBlockActor `mapstructure:",squash"`
}

// ID ...
func (*GlowItemFrame) ID() string {
	return IDGlowItemFrame
}
