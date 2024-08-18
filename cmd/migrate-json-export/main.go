package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/arran4/anytype-to-linkwarden"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	var LinkwardenToken = os.Getenv("LINKWARDEN_TOKEN")
	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	exportDir := flagSet.String("export-dir", "./Anytype.20240818.132330.55", "JSON Export directory to import")
	host := flagSet.String("linkwarden-endpoint", "https://linkwarden.com", "Linkwarden endpoint")
	dry := flagSet.Bool("dry", true, "Don't do anything (on by default)")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	exportDirContents, err := os.ReadDir(*exportDir)
	if err != nil {
		log.Fatal(err)
	}
	bookmarks := make([]*anytype_to_linkwarden.Bookmark, 0)
	tags := map[string]*anytype_to_linkwarden.AnyTypeTag{}
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
		ato := &anytype_to_linkwarden.AnytypeObject{}
		if err := json.Unmarshal(b, ato); err != nil {
			log.Fatal(err)
		}
		switch ato.SBType {
		case "Page":
			if slices.Contains(ato.Snapshot.Data.ObjectTypes, "ot-bookmark") {
				bookmarks = append(bookmarks, &anytype_to_linkwarden.Bookmark{ato})
			}
			fmt.Printf("Not a bookmark: %v\t%v\t%v\t%v\t%v\t%v\n", ato.SBType, ato.Snapshot.Data.Details.RelationKey, ato.Snapshot.Data.Details.Name, ato.Snapshot.Data.Details.Source, ato.Snapshot.Data.Details.Tag, ato.Snapshot.Data.ObjectTypes)
			continue
		case "STRelationOption":
			if ato.Snapshot.Data.Details.RelationKey == "tag" {
				tags[ato.Snapshot.Data.Details.Id] = &anytype_to_linkwarden.AnyTypeTag{ato}
			}
			fmt.Printf("Not a tag: %v\t%v\t%v\t%v\t%v\t%v\n", ato.SBType, ato.Snapshot.Data.Details.RelationKey, ato.Snapshot.Data.Details.Name, ato.Snapshot.Data.Details.Source, ato.Snapshot.Data.Details.Tag, ato.Snapshot.Data.ObjectTypes)
			continue
		default:
			fmt.Printf("Not something we know: %v\t%v\t%v\t%v\t%v\t%v\n", ato.SBType, ato.Snapshot.Data.Details.RelationKey, ato.Snapshot.Data.Details.Name, ato.Snapshot.Data.Details.Source, ato.Snapshot.Data.Details.Tag, ato.Snapshot.Data.ObjectTypes)
			continue
		}
	}
	slices.SortFunc(bookmarks, func(a, b *anytype_to_linkwarden.Bookmark) int {
		return strings.Compare(a.Snapshot.Data.Details.Name, b.Snapshot.Data.Details.Name)
	})

	if *dry {
		for _, bookmark := range bookmarks {
			fmt.Printf("Bookmark: %v\t%v\t%v\t%v\t%v\t%v\n", bookmark.SBType, bookmark.Snapshot.Data.Details.RelationKey, bookmark.Snapshot.Data.Details.Name, bookmark.Snapshot.Data.Details.Source, bookmark.Snapshot.Data.Details.Tags(tags), bookmark.Snapshot.Data.ObjectTypes)
		}
		return
	}

	collections, err := anytype_to_linkwarden.GetCollections(LinkwardenToken, *host)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Collections: %v", len(collections))
	for _, collection := range collections {
		log.Printf("Collection: %s %d %s", collection.Name, collection.Id, collection.Color)
	}

	collectionMap := maps.Collect(func(yield func(K string, V *anytype_to_linkwarden.Collection) bool) {
		for _, collection := range collections {
			yield(collection.Name, collection)
		}
	})

	anyTypeCollection, ok := collectionMap["AnyType"]
	if !ok {
		anyTypeCollection, err = anytype_to_linkwarden.CreateCollections(LinkwardenToken, *host, &anytype_to_linkwarden.PartialCreateCollection{
			Name:        "AnyType",
			Color:       "#0ea5e9",
			Description: "AnyType Imported data",
			IsPublic:    false,
		})
		if err != nil {
			log.Fatal(err)
		}
		collectionMap[anyTypeCollection.Name] = anyTypeCollection
	}
	log.Printf("CollectionMap: %#v", anyTypeCollection)

	for _, bookmark := range bookmarks {
		linkTagsIds := bookmark.Snapshot.Data.Details.Tags(tags)
		fmt.Printf("Bookmark: %v\t%v\t%v\t%v\t%v\t%v\n", bookmark.SBType, bookmark.Snapshot.Data.Details.RelationKey, bookmark.Snapshot.Data.Details.Name, bookmark.Snapshot.Data.Details.Source, linkTagsIds, bookmark.Snapshot.Data.ObjectTypes)
		var linkTags []*anytype_to_linkwarden.TagReference
		for _, tagStr := range linkTagsIds {
			linkTags = append(linkTags, &anytype_to_linkwarden.TagReference{
				Name: tagStr,
			})
		}
		newLink := &anytype_to_linkwarden.PartialCreateLink{
			Name: bookmark.Snapshot.Data.Details.Name,
			//Description: "Hi!",
			Url: bookmark.Snapshot.Data.Details.Source,
			Collection: &anytype_to_linkwarden.CollectionReference{
				Id:      &anyTypeCollection.Id,
				OwnerId: &anyTypeCollection.OwnerId,
				Name:    anyTypeCollection.Name,
			},
			Tags: linkTags,
		}

		link, err := anytype_to_linkwarden.PostLink(LinkwardenToken, *host, newLink)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Link: %v", link)
	}

}
