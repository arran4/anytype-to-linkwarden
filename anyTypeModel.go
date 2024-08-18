package anytype_to_linkwarden

import "slices"

type AnytypeObjectDetails struct {
	Tag         []string `json:"tag,omitempty"`
	Source      string   `json:"source,omitempty"`
	Name        string   `json:"name,omitempty"`
	Id          string   `json:"id,omitempty"`
	RelationKey string   `json:"relationKey,omitempty"`
}

func (d *AnytypeObjectDetails) Tags(tags map[string]*AnyTypeTag) []string {
	return slices.Collect(func(yield func(E string) bool) {
		for _, tag := range d.Tag {
			v := tags[tag]
			if v != nil {
				if !yield(v.Snapshot.Data.Details.Name) {
					return
				}
			}
		}
	})
}

type AnytypeObjectData struct {
	ObjectTypes []string             `json:"objectTypes,omitempty"`
	Details     AnytypeObjectDetails `json:"details,omitempty"`
}

type AnytypeObjectSnapshot struct {
	Data AnytypeObjectData `json:"data,omitempty"`
}

type AnytypeObject struct {
	SBType   string                `json:"sbType,omitempty"`
	Snapshot AnytypeObjectSnapshot `json:"snapshot,omitempty"`
}

type Bookmark struct {
	*AnytypeObject
}

type AnyTypeTag struct {
	*AnytypeObject
}
