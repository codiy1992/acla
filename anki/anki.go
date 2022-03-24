package anki

import (
	"log"
	"regexp"
	"strings"

	"github.com/codiy1992/anki-connect-go"
)

var client anki.Client

func init() {
	client = *anki.NewDefaultClient()
}

func Client() *anki.Client {
	return &client
}

func IsMainField(model string, fieldName string) bool {
	fields, err := client.ModelFieldNames(model)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}
	if len(fields) == 0 {
		log.Fatalf("Model `%s` does not exist", model)
	}
	if fieldName != fields[0] {
		log.Printf("filed `%s` does match first field `%s` of model `%s`", fieldName, fields[0], model)
		return false
	}
	return true
}

func AddNote(deck string, model string, fields map[string]string, tags []string) int {

	for field, value := range fields {
		if strings.HasSuffix(field, "s") {
			re := regexp.MustCompile(`<\s*img[^>]*>`)
			if !re.MatchString(value) {
				fields[field] = "<ul><li>" + value + "</li></ul>"
			}
		}
	}
	note := anki.NoteInput{
		Deck:   deck,
		Model:  model,
		Fields: fields,
		Options: map[string]interface{}{
			"allowDuplicate": false,
			"duplicateScope": "deck",
			"duplicateScopeOption": map[string]interface{}{
				"deckName":       deck,
				"checkChildren":  false,
				"checkAllModels": false,
			},
		},
		Tags: tags,
	}
	noteId, err := client.AddNote(note)
	if err != nil {
		log.Fatalf("Create note failed: %s", err)
	}
	return noteId
}

func UpdateNote(query string, fields map[string]string, tags []string, override bool) {

	notes, _ := client.FindNotes(query)

	if len(notes) == 0 {
		log.Fatalf("No notes found!")
	}
	if len(notes) > 1 {
		log.Fatalf("More than one note found!")
	}

	if !override {
		res2, _ := client.NotesInfo(notes)
		for field, value := range fields {
			if res2[0].Fields[field] == "" {
				log.Fatalf("field `%s` doesn't exists", field)
			}
			original := res2[0].Fields[field].(map[string]interface{})
			fields[field] = strings.TrimRight(original["value"].(string), "\n") + "\n" + value
			if strings.Trim(original["value"].(string), "\n ") == "" {
				re := regexp.MustCompile(`<\s*img[^>]*>`)
				if !re.MatchString(value) {
					fields[field] = "<ul><li>" + value + "</li></ul>"
				}
			} else {
				re := regexp.MustCompile(`</li>[^<]*</ul>`)
				if re.MatchString(original["value"].(string)) {
					fields[field] = re.ReplaceAllString(
						original["value"].(string),
						"</li><li>"+value+"</li></ul>",
					)
				}
			}
		}
	}

	note := anki.NoteFieldsInput{
		ID:     notes[0],
		Fields: fields,
	}

	if err := client.UpdateNote(note); err != nil {
		log.Println("update failed:", err)
	}

	// TODO addTags
}

func SuspendCard(query string) bool {
	cardIds, err := client.FindCards(query)
	if err != nil {
		log.Fatalf("FindCards failed with query string `%s`", query)
	}
	if len(cardIds) > 1 {
		log.Fatalf("Found more than one cards, please check your query string `%s`", query)
	}
	suspended, err := client.Suspended(cardIds[0])
	if err != nil {
		log.Fatalf("Cloud not check card %d suspended or not", cardIds[0])
	}
	if suspended {
		return true
	}
	result, err := client.Suspend(cardIds)
	if err != nil {
		log.Fatalf("Cards %v suspend failed: %s", cardIds, err)
	}
	return result
}
