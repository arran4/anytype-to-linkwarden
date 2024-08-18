package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	LinkwardenToken = os.Getenv("LINKWARDEN_TOKEN")
)

type AnytypeObjectDetails struct {
	Tag         []string `json:"tag,omitempty"`
	Source      string   `json:"source,omitempty"`
	Name        string   `json:"name,omitempty"`
	Id          string   `json:"id,omitempty"`
	RelationKey string   `json:"relationKey,omitempty"`
}

func (d *AnytypeObjectDetails) Tags(tags map[string]*Tag) []string {
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

type Tag struct {
	*AnytypeObject
}

func main() {
	flags := flag.NewFlagSet("migrate-json-export", flag.ExitOnError)
	exportDir := flags.String("export-dir", "./Anytype.20240818.132330.55", "JSON Export directory to import")
	err := flags.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	exportDirContents, err := os.ReadDir(*exportDir)
	if err != nil {
		log.Fatal(err)
	}
	bookmarks := make([]*Bookmark, 0)
	tags := map[string]*Tag{}
	for _, exportDirContent := range exportDirContents {
		if exportDirContent.IsDir() {
			continue
		}
		if !strings.HasSuffix(exportDirContent.Name(), ".pb.json") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(*exportDir, exportDirContent.Name()))
		if err != nil {
			log.Fatal(err)
		}
		ato := &AnytypeObject{}
		if err := json.Unmarshal(b, ato); err != nil {
			log.Fatal(err)
		}
		switch ato.SBType {
		case "Page":
			if slices.Contains(ato.Snapshot.Data.ObjectTypes, "ot-bookmark") {
				bookmarks = append(bookmarks, &Bookmark{ato})
			}
			fmt.Printf("Not a bookmark: %v\t%v\t%v\t%v\t%v\t%v\n", ato.SBType, ato.Snapshot.Data.Details.RelationKey, ato.Snapshot.Data.Details.Name, ato.Snapshot.Data.Details.Source, ato.Snapshot.Data.Details.Tag, ato.Snapshot.Data.ObjectTypes)
			continue
		case "STRelationOption":
			if ato.Snapshot.Data.Details.RelationKey == "tag" {
				tags[ato.Snapshot.Data.Details.Id] = &Tag{ato}
			}
			fmt.Printf("Not a tag: %v\t%v\t%v\t%v\t%v\t%v\n", ato.SBType, ato.Snapshot.Data.Details.RelationKey, ato.Snapshot.Data.Details.Name, ato.Snapshot.Data.Details.Source, ato.Snapshot.Data.Details.Tag, ato.Snapshot.Data.ObjectTypes)
			continue
		default:
			fmt.Printf("Not something we know: %v\t%v\t%v\t%v\t%v\t%v\n", ato.SBType, ato.Snapshot.Data.Details.RelationKey, ato.Snapshot.Data.Details.Name, ato.Snapshot.Data.Details.Source, ato.Snapshot.Data.Details.Tag, ato.Snapshot.Data.ObjectTypes)
			continue
		}
	}
	slices.SortFunc(bookmarks, func(a, b *Bookmark) int {
		return strings.Compare(a.Snapshot.Data.Details.Name, b.Snapshot.Data.Details.Name)
	})
	for _, bookmark := range bookmarks {
		fmt.Printf("Bookmark: %v\t%v\t%v\t%v\t%v\t%v\n", bookmark.SBType, bookmark.Snapshot.Data.Details.RelationKey, bookmark.Snapshot.Data.Details.Name, bookmark.Snapshot.Data.Details.Source, bookmark.Snapshot.Data.Details.Tags(tags), bookmark.Snapshot.Data.ObjectTypes)
	}

}
