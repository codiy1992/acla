package dict

func Query(source string, word string) interface{} {
	if source == "cambridge" {
		return NewCambridge().Query(word)
	} else {
		return NewCambridge().Query(word)
	}
}
