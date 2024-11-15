package skin_process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

func ProcessGeometry(skin *Skin, rawData []byte) (err error) {
	/* Layer 1 */
	geometryMap := map[string]json.RawMessage{}
	if err = json.Unmarshal(rawData, &geometryMap); err != nil {
		return fmt.Errorf("ProcessGeometry: %v", err)
	}
	// setup resource patch and get geometry data
	var skinGeometry json.RawMessage
	var geometryName string
	var skinResourcePatchReplaceString string
	for k, v := range geometryMap {
		if strings.HasPrefix(k, "geometry.") {
			geometryName = k
			skinGeometry = v
			break
		}
	}
	if geometryName == "" {
		return fmt.Errorf("ProcessGeometry: lack of geometry data")
	}
	if skin.SkinIsSlim {
		skinResourcePatchReplaceString = "geometry.humanoid.customSlim"
	} else {
		skinResourcePatchReplaceString = "geometry.humanoid.custom"
	}
	skin.SkinResourcePatch = bytes.ReplaceAll(
		skin.SkinResourcePatch,
		[]byte(skinResourcePatchReplaceString),
		[]byte(geometryName),
	)
	/* Layer 2 */
	geometry := &SkinGeometry{}
	if err = json.Unmarshal(skinGeometry, geometry); err != nil {
		return fmt.Errorf("ProcessGeometry: %v", err)
	}
	// handle bones
	hasRoot := false
	renderGroupNames := []string{"leftArm", "rightArm"}
	for _, bone := range geometry.Bones {
		// setup parent
		switch bone.Name {
		case "waist", "leftLeg", "rightLeg":
			bone.Parent = "root"
		case "head":
			bone.Parent = "body"
		case "leftArm", "rightArm":
			bone.Parent = "body"
			bone.RenderGroupID = 1
		case "body":
			bone.Parent = "waist"
		case "root":
			hasRoot = true
		}
		// setup render group
		if slices.Contains(renderGroupNames, bone.Parent) {
			bone.RenderGroupID = 1
			renderGroupNames = append(renderGroupNames, bone.Name)
		}
	}
	if !hasRoot {
		geometry.Bones = append(geometry.Bones, &SkinGeometryBone{
			Name: "root",
			Pivot: []json.Number{
				json.Number("0.0"),
				json.Number("0.0"),
				json.Number("0.0"),
			},
		})
	}
	// return
	skin.SkinGeometry, _ = json.Marshal(map[string]any{
		"format_version": "1.8.0",
		geometryName:     geometry,
	})
	return
}
