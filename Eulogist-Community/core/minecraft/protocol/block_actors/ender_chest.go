package block_actors

import general "Eulogist/core/minecraft/protocol/block_actors/general_actors"

// 末影箱
type EnderChest struct {
	general.ChestBlockActor `mapstructure:",squash"`
}

// ID ...
func (*EnderChest) ID() string {
	return IDEnderChest
}
