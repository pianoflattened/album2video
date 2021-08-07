package main

import (
	"strings"
	"unicode"

	whatlang "github.com/abadojack/whatlanggo"
)

var nocapEng = []string{"a", "an", "the", "and", "but", "or", "nor", "for", "yet", "so", "versus", "v.", "vs.", "etc.", "et", "cetera", "as", "by", "in", "on", "to", "n'", "o'"}
var nocapTurk = []string{"ve", "ile", "ya", "veya", "yahut", "ki", "da", "de"}
var qPart = []string{"mı", "mi", "mu", "mü"}

// could be more concise with gotos
func TitleCase(s string, selfcall bool) string {
	words := strings.Fields(s)
	var first, rest string
	for i, word := range words {
		if len(word) == 1 {
			if selfcall {
				words[i] = strings.ToTitle(word)
				continue
			}

			script := whatlang.DetectScript(word)
			switch script {
			case unicode.Latin:
				lang := whatlang.DetectLang(strings.Join(words, " "))
				if lang == whatlang.Eng {
					if strings.ToLower(word) == "a" {
						words[i] = "a"
					} else {
						words[i] = strings.ToTitle(words[i])
					}
					continue
				} else if lang == whatlang.Tur {
					words[i] = strings.ToTitle(words[i])
					continue
				}
				fallthrough
			case unicode.Cyrillic:
				fallthrough
			case unicode.Greek:
				fallthrough
			case unicode.Armenian:
				words[i] = strings.ToTitle(words[i])
			}
			continue
		}

		first = string([]byte{word[0]})
		rest = word[1:]

		script := whatlang.DetectScript(word)
		switch script {
		case unicode.Latin:
			lang := whatlang.DetectLang(word)
			if lang == whatlang.Eng {
				if i == 0 && !selfcall {
					words[i] = strings.ToTitle(first) + strings.ToLower(rest)
					continue
				}

				isException := false
				for _, exception := range nocapEng {
					if strings.ToLower(word) == exception {
						isException = true
						break
					}
				}
				if isException {
					words[i] = strings.ToLower(first + rest)
					continue
				}

				if len(strings.Split(word, ".")) >= len(strings.ReplaceAll(word, ".", "")) { // acronym with dots like V.O.S. Family
					words[i] = strings.ToTitle(first + rest)
					continue
				}

				// this one hurts my brain to think about and will probably bug out
				if strings.ContainsRune(word, '-') {
					words[i] = strings.ReplaceAll(TitleCase(strings.ReplaceAll(word, "-", " "), true), " ", "-")
				}

			} else {
				switch lang {
				case whatlang.Tur:
					if i == 0 && !selfcall {
						words[i] = strings.ToTitle(first) + strings.ToLower(rest)
						continue
					}

					if len(word) >= 2 {
						isQPart := false
						for _, prefix := range qPart {
							if strings.HasPrefix(strings.ToLower(word), prefix) {
								isQPart = true
								break
							}
						}

						if isQPart {
							words[i] = strings.ToLower(first + rest)
						}
						continue
					}

					wordInNoCap := false
					for _, conj := range nocapTurk {
						if conj == word {
							wordInNoCap = true
							break
						}
					}

					if wordInNoCap {
						words[i] = strings.ToLower(first + rest)
						continue
					}

					words[i] = strings.ToTitle(first) + strings.ToLower(rest)
					continue

				case whatlang.Nld:
					if len(word) > 2 && i == 0 && !selfcall {
						if strings.ToLower(word[0:2]) == "ij" {
							first = "ij"
							rest = word[2:]
						}
					}
					fallthrough
				default:
					if i == 0 && !selfcall { /* pass */
					} else {
						words[i] = strings.ToLower(first + rest)
						continue
					}
					words[i] = strings.ToTitle(first) + strings.ToLower(rest)
				}
			}
		case unicode.Cyrillic: // these have the same rules (1st letter of 1st word, proper nouns)
			fallthrough // i am not doing named entity extraction or anything like that
		case unicode.Greek: // i'll put in a keyboard shortcut that replaces selection with its
			fallthrough // corresponding capital letter
		case unicode.Armenian:
			if i == 0 && !selfcall {
				first = strings.ToTitle(first)
			} else {
				first = strings.ToLower(first)
			}
			words[i] = first + strings.ToLower(rest)
		default:
			// no capitalization required
			continue
		}
	}

	return strings.Join(words, " ")
}
