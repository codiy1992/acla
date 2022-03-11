package main

import (
	"context"
	"flag"
	"log"

	"github.com/leonhfr/anki-connect-go"
)

var flagDeckName = flag.String("deck", "current", "deck name")

func main() {
	flag.Parse()
	client := anki.NewDefaultClient()
	notes, _ := client.FindNotes(context.TODO(), "deck:English Word:prevent")
	res2, _ := client.NotesInfo(context.TODO(), notes)
	log.Println(res2)
	// if deckName := *flagDeckName; deckName != "" {
	// log.Printf("Please specify the deck name!")
	// return
	// }
}
