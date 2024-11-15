package block_actors

import general "Eulogist/core/minecraft/protocol/block_actors/general_actors"

// 高炉
type BlastFurnace struct {
	general.FurnaceBlockActor `mapstructure:",squash"`
}

// ID ...
func (*BlastFurnace) ID() string {
	return IDBlastFurnace
}
