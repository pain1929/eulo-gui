package block_actors

import general "Eulogist/core/minecraft/protocol/block_actors/general_actors"

// 阳光探测器
type DayLightDetector struct {
	general.BlockActor `mapstructure:",squash"`
}

// ID ...
func (*DayLightDetector) ID() string {
	return IDDayLightDetector
}
