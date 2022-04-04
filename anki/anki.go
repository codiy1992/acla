package anki

import "log"

func SuspendCard(query string) bool {
	cardIds, err := client.FindCards(query)
	if err != nil {
		log.Fatalf("FindCards failed with query string `%s`", query)
	}
	if len(cardIds) > 1 {
		log.Fatalf("Found more than one cards, please check your query string `%s`", query)
	}
	if len(cardIds) == 0 {
		log.Fatalf("No cards found, please check your query string `%s`", query)
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
