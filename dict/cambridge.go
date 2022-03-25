package dict

import (
	"encoding/json"
	"fmt"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Vocabulary struct {
	Dictionary    string                   `json:"dictionary"`
	Pronunciation map[string]Pronunciation `json:"pronunciation"`
	Definitions   []Definition             `json:"definitions"`
	Examples      []string                 `json:"examples"`
}

type Pronunciation struct {
	Ipa string `json:"ipa"`
	Mp3 string `json:"mp3"`
	Ogg string `json:"ogg"`
}

type Definition struct {
	Definition string   `json:"definition"`
	Images     []string `json:"images"`
	Examples   []string `json:"examples"`
	Thesaurus  []string `json:"thesaurus"`
}

func Query(vocab string) []Vocabulary {
	doc, err := htmlquery.LoadURL("https://dictionary.cambridge.org/us/dictionary/english/" + vocab)
	if err != nil {
		panic(err)
	}
	dictionaries, err := htmlquery.QueryAll(doc, `//div[@class="pr dictionary"]`)
	if err != nil {
		panic(err)
	}
	var result []Vocabulary
	for _, dictionary := range dictionaries {
		Vocabulary := Vocabulary{}
		Vocabulary.Dictionary = "Advanced Learner's Dictionary"
		dictName, _ := htmlquery.Query(dictionary, `//div[@class="di-title di_t"]/h2[contains(., "Intermediate")]`)
		if dictName != nil {
			Vocabulary.Dictionary = "Intermediate English Dictionary"
			fmt.Println("--------------------> Intermediate English Dictionary")
		} else {
			dictName, _ := htmlquery.Query(dictionary, `//div[@class="di-title di_t"]/h2[contains(., "Business")]`)
			if dictName != nil {
				Vocabulary.Dictionary = "Business English Dictionary"
				fmt.Println("--------------------> Business English Dictionary")
			}
		}

		entry_body, _ := htmlquery.Query(dictionary, `//div[@class="entry-body"]`)
		fmt.Println("---------- entry body ----------")
		Vocabulary.Pronunciation = ParsePronuns(entry_body)
		// definitions
		fmt.Println("----- definitions -----")
		definitions, _ := htmlquery.QueryAll(entry_body, `//div[@class="def-block ddef_block "]`)
		for _, definition := range definitions {
			Vocabulary.Definitions = append(Vocabulary.Definitions, ParseDefinition(definition))
		}
		fmt.Println("----- definitions -----")
		// more examples
		examples, _ := htmlquery.QueryAll(entry_body, `//div[contains(@class, "sense-body")]/div[@class="daccord"]//section/ul/li[@class="eg dexamp hax"]`)
		for _, example := range examples {
			Vocabulary.Examples = append(Vocabulary.Examples, htmlquery.InnerText(example))
		}
		fmt.Println("---------- entry body ----------")

		result = append(result, Vocabulary)
	}
	json, _ := json.Marshal(result)
	fmt.Println(string(json))
	return result
}

func ParsePronuns(node *html.Node) map[string]Pronunciation {
	result := make(map[string]Pronunciation)

	ipa_us, _ := htmlquery.Query(node, `//span[@class="us dpron-i "][1]//span[contains(@class, "ipa")]`)
	if ipa_us != nil {
		Pronunciation := Pronunciation{}
		Pronunciation.Ipa = htmlquery.InnerText(ipa_us)
		fmt.Println(htmlquery.InnerText(ipa_us))

		audio_mp3, _ := htmlquery.Query(node, `//span[@class="us dpron-i "][1]//source[@type="audio/mpeg"]`)
		if audio_mp3 != nil {
			Pronunciation.Mp3 = htmlquery.SelectAttr(audio_mp3, "src")
			fmt.Println(htmlquery.SelectAttr(audio_mp3, "src"))
		}
		audio_ogg, _ := htmlquery.Query(node, `//span[@class="us dpron-i "][1]//source[@type="audio/ogg"]`)
		if audio_ogg != nil {
			Pronunciation.Ogg = htmlquery.SelectAttr(audio_ogg, "src")
			fmt.Println(htmlquery.SelectAttr(audio_ogg, "src"))
		}
		result["us"] = Pronunciation
	}

	ipa_uk, _ := htmlquery.Query(node, `//span[@class="uk dpron-i "][1]//span[contains(@class, "ipa")]`)
	if ipa_uk != nil {
		Pronunciation := Pronunciation{}
		Pronunciation.Ipa = htmlquery.InnerText(ipa_us)
		fmt.Println(htmlquery.InnerText(ipa_uk))
		audio_mp3, _ := htmlquery.Query(node, `//span[@class="uk dpron-i "][1]//source[@type="audio/mpeg"]`)
		if audio_mp3 != nil {
			Pronunciation.Mp3 = htmlquery.SelectAttr(audio_mp3, "src")
			fmt.Println(htmlquery.SelectAttr(audio_mp3, "src"))
		}
		audio_ogg, _ := htmlquery.Query(node, `//span[@class="uk dpron-i "][1]//source[@type="audio/ogg"]`)
		if audio_ogg != nil {
			Pronunciation.Ogg = htmlquery.SelectAttr(audio_ogg, "src")
			fmt.Println(htmlquery.SelectAttr(audio_ogg, "src"))
		}
		result["uk"] = Pronunciation
	}
	return result
}

func ParseDefinition(node *html.Node) Definition {
	Definition := Definition{}
	define, _ := htmlquery.Query(node, `//div[@class="def ddef_d db"]`)
	if define != nil {
		Definition.Definition = htmlquery.InnerText(define)
		fmt.Println(htmlquery.InnerText(define))
	}

	// images
	images, _ := htmlquery.QueryAll(node, `//div[@class="dimg"]//amp-img`)
	for _, image := range images {
		Definition.Images = append(Definition.Images, htmlquery.SelectAttr(image, "on"))
		fmt.Println(htmlquery.SelectAttr(image, "on"))
	}

	// examples
	examples, _ := htmlquery.QueryAll(node, `//div[@class="examp dexamp"]`)
	for _, example := range examples {
		Definition.Examples = append(Definition.Examples, htmlquery.InnerText(example))
		fmt.Println(htmlquery.InnerText(example))
	}

	// thesaurus
	thesaurus, _ := htmlquery.Query(node, `//section//header[contains(., "Thesaurus")]`)
	if thesaurus != nil {
		Definition.Thesaurus = append(Definition.Thesaurus, htmlquery.InnerText(thesaurus))
		fmt.Println(htmlquery.InnerText(thesaurus))
	}
	return Definition
}
