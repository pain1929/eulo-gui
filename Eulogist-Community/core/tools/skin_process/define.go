package skin_process

import (
	_ "embed"
	"encoding/json"
)

// ----------------------------------------------------------------------------------------------------

//go:embed default_wide_skin_resource_patch.json
var DefaultWideSkinResourcePatch []byte

//go:embed default_slim_skin_resource_patch.json
var DefaultSlimSkinResourcePatch []byte

//go:embed default_skin_geometry.json
var DefaultSkinGeometry []byte

//go:embed steve.png
var SteveSkin []byte

// ----------------------------------------------------------------------------------------------------

// 描述皮肤信息
type Skin struct {
	// 储存皮肤数据的二进制负载。
	// 对于普通皮肤，这是一个二进制形式的 PNG；
	// 对于高级皮肤(如 4D 皮肤)，
	// 这是一个压缩包形式的二进制表示
	FullSkinData []byte
	// 皮肤的 UUID
	SkinUUID string
	// 皮肤项目的 UUID
	SkinItemID string
	// 皮肤的手臂是否纤细
	SkinIsSlim bool
	// 皮肤的一维密集像素矩阵
	SkinPixels []byte
	// 皮肤的骨架信息
	SkinGeometry []byte
	// SkinResourcePatch 是一个 JSON 编码对象，
	// 其中包含一些指向皮肤所具有的几何形状的字段。
	// 它包含的 JSON 对象指定动画的几何形状，
	// 以及播放器的默认皮肤的组合方式
	SkinResourcePatch []byte
	// 皮肤的宽度
	SkinWidth int
	// 皮肤的高度
	SkinHight int
}

// ----------------------------------------------------------------------------------------------------

type SkinCube struct {
	Inflate *json.Number  `json:"inflate,omitempty"`
	Mirror  *bool         `json:"mirror,omitempty"`
	Origin  []json.Number `json:"origin"`
	Size    []json.Number `json:"size"`
	Uv      []json.Number `json:"uv"`
}

type SkinGeometryBone struct {
	Cubes         *[]SkinCube   `json:"cubes,omitempty"`
	Name          string        `json:"name"`
	Parent        string        `json:"parent,omitempty"`
	Pivot         []json.Number `json:"pivot"`
	RenderGroupID int           `json:"render_group_id,omitempty"`
	Rotation      []json.Number `json:"rotation,omitempty"`
}

type SkinGeometry struct {
	Bones               []*SkinGeometryBone `json:"bones"`
	TextureHeight       int                 `json:"textureheight"`
	TextureWidth        int                 `json:"texturewidth"`
	VisibleBoundsHeight json.Number         `json:"visible_bounds_height,omitempty"`
	VisibleBoundsOffset []json.Number       `json:"visible_bounds_offset,omitempty"`
	VisibleBoundsWidth  json.Number         `json:"visible_bounds_width,omitempty"`
}

// ----------------------------------------------------------------------------------------------------

// ...
type SkinManifest struct {
	Header struct {
		UUID string `json:"uuid"`
	} `json:"header"`
}
