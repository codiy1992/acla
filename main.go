package main

import (
	"acla/internal"
	"flag"
	"strings"
)

var flagDeckName = flag.String("deck", "current", "deck name")
var flagModelName = flag.String("model", "", "model name(note type)")
var flagFieldName = flag.String("field", "", "filed name")
var flagFieldValue = flag.String("value", "", "filed value")
var flagQuery = flag.String("query", "", "query string")
var flagTags = flag.String("tags", "", "tags")

func main() {
	flag.Parse()

	if *flagQuery == "" {
		internal.CreateNote(
			*flagDeckName,
			*flagModelName,
			*flagFieldName,
			*flagFieldValue,
			strings.Split(*flagTags, ","),
		)
	} else {
		internal.UpdateFiled(*flagQuery, *flagFieldName, *flagFieldValue, true)
	}
}
