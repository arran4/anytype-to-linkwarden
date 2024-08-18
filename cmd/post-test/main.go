package main

import (
	"flag"
	"github.com/arran4/anytype-to-linkwarden"
	"log"
	"maps"
	"os"
)

func main() {
	var LinkwardenToken = os.Getenv("LINKWARDEN_TOKEN")
	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	host := flagSet.String("linkwarden-endpoint", "https://linkwarden.com", "Linkwarden endpoint")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	//tags, err := anytype_to_linkwarden.GetTags(LinkwardenToken, *host)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("Tags: %v", len(tags))
	//for _, tag := range tags {
	//	log.Printf("Tag: %s %d", tag.Name, tag.Id)
	//}

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

	//tagMap := maps.Collect(func(yield func(K string, V *anytype_to_linkwarden.Tag) bool) {
	//	for _, tag := range tags {
	//		yield(tag.Name, tag)
	//	}
	//})

	newLink := &anytype_to_linkwarden.PartialCreateLink{
		//Name:        "Hi!",
		//Description: "Hi!",
		Url: "https://www.arran.net.au/",
		Collection: &anytype_to_linkwarden.CollectionReference{
			Id:      &anyTypeCollection.Id,
			OwnerId: &anyTypeCollection.OwnerId,
			Name:    anyTypeCollection.Name,
		},
		Tags: []*anytype_to_linkwarden.TagReference{
			{
				Name: "Test Tag 1",
			},
			{
				Name: "Test Tag 2",
			},
		},
	}

	link, err := anytype_to_linkwarden.PostLink(LinkwardenToken, *host, newLink)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Link: %v", link)

}
