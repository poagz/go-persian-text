package persian

import (
	"bytes"
	"strings"
	"unicode"
)

// letterGroup represents the letter and bounding letters
type letterGroup struct {
	backLetter  rune
	letter      rune
	frontLetter rune
}

// letterShape represents all shapes of persian letters in a word
type letterShape struct {
	Isolated, Initial, Medial, Final rune
}

// Map of different shapes of persian alphabet
var persianAlphabetShapes = map[rune]letterShape{
	// Letter (ﺁ)
	'\u0622': {Isolated: '\uFE81', Initial: '\u0622', Medial: '\uFE82', Final: '\uFE82'},
	// Letter (ﺍ)
	'\u0627': {Isolated: '\uFE8D', Initial: '\u0627', Medial: '\uFE8E', Final: '\uFE8E'},
	// Letter (ﺏ)
	'\u0628': {Isolated: '\uFE8F', Initial: '\uFE91', Medial: '\uFE92', Final: '\uFE90'},
	// Letter (ﭖ)
	'\u067E': {Isolated: '\uFB56', Initial: '\uFB58', Medial: '\uFB59', Final: '\uFB57'},
	// Letter (ﺕ)
	'\u062A': {Isolated: '\uFE95', Initial: '\uFE97', Medial: '\uFE98', Final: '\uFE96'},
	// Letter (ﺙ)
	'\u062B': {Isolated: '\uFE99', Initial: '\uFE9B', Medial: '\uFE9C', Final: '\uFE9A'},
	// Letter (ﺝ)
	'\u062C': {Isolated: '\uFE9D', Initial: '\uFE9F', Medial: '\uFEA0', Final: '\uFE9E'},
	// Letter (ﭺ)
	'\u0686': {Isolated: '\uFB7A', Initial: '\uFB7C', Medial: '\uFB7D', Final: '\uFB7B'},
	// Letter (ﺡ)
	'\u062D': {Isolated: '\uFEA1', Initial: '\uFEA3', Medial: '\uFEA4', Final: '\uFEA2'},
	// Letter (ﺥ)
	'\u062E': {Isolated: '\uFEA5', Initial: '\uFEA7', Medial: '\uFEA8', Final: '\uFEA6'},
	// Letter (ﺩ)
	'\u062F': {Isolated: '\uFEA9', Initial: '\u062F', Medial: '\uFEAA', Final: '\uFEAA'},
	// Letter (ﺫ)
	'\u0630': {Isolated: '\uFEAB', Initial: '\u0630', Medial: '\uFEAC', Final: '\uFEAC'},
	// Letter (ﺭ)
	'\u0631': {Isolated: '\uFEAD', Initial: '\u0631', Medial: '\uFEAE', Final: '\uFEAE'},
	// Letter (ﺯ)
	'\u0632': {Isolated: '\uFEAF', Initial: '\u0632', Medial: '\uFEB0', Final: '\uFEB0'},
	// Letter (ﮊ)
	'\u0698': {Isolated: '\uFB8A', Initial: '\u0698', Medial: '\uFB8B', Final: '\uFB8B'},
	// Letter (ﺱ)
	'\u0633': {Isolated: '\uFEB1', Initial: '\uFEB3', Medial: '\uFEB4', Final: '\uFEB2'},
	// Letter (ﺵ)
	'\u0634': {Isolated: '\uFEB5', Initial: '\uFEB7', Medial: '\uFEB8', Final: '\uFEB6'},
	// Letter (ﺹ)
	'\u0635': {Isolated: '\uFEB9', Initial: '\uFEBB', Medial: '\uFEBC', Final: '\uFEBA'},
	// Letter (ﺽ)
	'\u0636': {Isolated: '\uFEBD', Initial: '\uFEBF', Medial: '\uFEC0', Final: '\uFEBE'},
	// Letter (ﻁ)
	'\u0637': {Isolated: '\uFEC1', Initial: '\uFEC3', Medial: '\uFEC4', Final: '\uFEC2'},
	// Letter (ﻅ)
	'\u0638': {Isolated: '\uFEC5', Initial: '\uFEC7', Medial: '\uFEC8', Final: '\uFEC6'},
	// Letter (ﻉ)
	'\u0639': {Isolated: '\uFEC9', Initial: '\uFECB', Medial: '\uFECC', Final: '\uFECA'},
	// Letter (ﻍ)
	'\u063A': {Isolated: '\uFECD', Initial: '\uFECF', Medial: '\uFED0', Final: '\uFECE'},
	// Letter (ﻑ)
	'\u0641': {Isolated: '\uFED1', Initial: '\uFED3', Medial: '\uFED4', Final: '\uFED2'},
	// Letter (ﻕ)
	'\u0642': {Isolated: '\uFED5', Initial: '\uFED7', Medial: '\uFED8', Final: '\uFED6'},
	// Letter (ﮎ)
	'\u06A9': {Isolated: '\uFB8E', Initial: '\uFB90', Medial: '\uFB91', Final: '\uFB8F'},
	// Letter (ﮒ)
	'\u06AF': {Isolated: '\uFB92', Initial: '\uFB94', Medial: '\uFB95', Final: '\uFB93'},
	// Letter (ﻝ)
	'\u0644': {Isolated: '\uFEDD', Initial: '\uFEDF', Medial: '\uFEE0', Final: '\uFEDE'},
	// Letter (ﻡ)
	'\u0645': {Isolated: '\uFEE1', Initial: '\uFEE3', Medial: '\uFEE4', Final: '\uFEE2'},
	// Letter (ﻥ)
	'\u0646': {Isolated: '\uFEE5', Initial: '\uFEE7', Medial: '\uFEE8', Final: '\uFEE6'},
	// Letter (ﻭ)
	'\u0648': {Isolated: '\uFEED', Initial: '\u0648', Medial: '\uFEEE', Final: '\uFEEE'},
	// Letter (ﻩ)
	'\u0647': {Isolated: '\uFEE9', Initial: '\uFEEB', Medial: '\uFEEC', Final: '\uFEEA'},
	// Letter (ی)
	'\u06CC': {Isolated: '\uFBFC', Initial: '\uFBFE', Medial: '\uFBFF', Final: '\uFBFD'},
}

// ReShape will reconstruct persian text to be connected correctly
func ReShape(input string) string {
	var langSections []string
	var continuousLangFa string
	var continuousLangLt string

	for _, letter := range input {
		if IsPersianLetter(letter) {
			if len(continuousLangLt) > 0 {
				langSections = append(langSections, strings.TrimSpace(continuousLangLt))
			}
			continuousLangLt = ""
			continuousLangFa += string(letter)
		} else {
			if len(continuousLangFa) > 0 {
				langSections = append(langSections, strings.TrimSpace(continuousLangFa))
			}
			continuousLangFa = ""
			continuousLangLt += string(letter)
		}
	}
	if len(continuousLangLt) > 0 {
		langSections = append(langSections, strings.TrimSpace(continuousLangLt))
	}
	if len(continuousLangFa) > 0 {
		langSections = append(langSections, strings.TrimSpace(continuousLangFa))
	}

	var shapedSentence []string
	for _, section := range langSections {
		if IsPersian(section) {
			for _, word := range strings.Fields(section) {
				shapedSentence = append(shapedSentence, shapeWord(word))
			}
		} else {
			shapedSentence = append(shapedSentence, section)
		}
	}

	//Reverse words
	for i := len(shapedSentence)/2 - 1; i >= 0; i-- {
		opp := len(shapedSentence) - 1 - i
		shapedSentence[i], shapedSentence[opp] = shapedSentence[opp], shapedSentence[i]
	}
	return strings.Join(shapedSentence, " ")
}

// shapeWord will reconstruct a persian word to be connected correctly
func shapeWord(input string) string {
	if !IsPersian(input) {
		return input
	}

	var shapedInput bytes.Buffer

	//Convert input into runes
	inputRunes := []rune(input)
	for i := range inputRunes {
		//Get Bounding back and front letters
		var backLetter, frontLetter rune
		if i-1 >= 0 {
			backLetter = inputRunes[i-1]
		}
		if i != len(inputRunes)-1 {
			frontLetter = inputRunes[i+1]
		}
		//Fix the letter based on bounding letters
		if _, ok := persianAlphabetShapes[inputRunes[i]]; ok {
			adjustedLetter := adjustLetter(letterGroup{backLetter, inputRunes[i], frontLetter})
			shapedInput.WriteRune(adjustedLetter)
		} else {
			shapedInput.WriteRune(inputRunes[i])
		}
	}

	return reverse(shapedInput.String())
}

// reverse the persian string for RTL support in rendering
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// adjustLetter will adjust the persian letter depending on its position
func adjustLetter(g letterGroup) rune {

	switch {
	//In between 2 letters
	case g.backLetter > 0 && g.frontLetter > 0:
		if isNoneLeftJoiningLetter(g.backLetter) {
			return persianAlphabetShapes[g.letter].Initial
		}
		return persianAlphabetShapes[g.letter].Medial

	//Not preceded by any letter
	case g.backLetter == 0 && g.frontLetter > 0:
		return persianAlphabetShapes[g.letter].Initial

	//Not followed by any letter
	case g.backLetter > 0 && g.frontLetter == 0:
		if isNoneLeftJoiningLetter(g.backLetter) {
			return persianAlphabetShapes[g.letter].Isolated
		}
		return persianAlphabetShapes[g.letter].Final

	default:
		return persianAlphabetShapes[g.letter].Isolated
	}
}

// letters become after these group as initial form
func isNoneLeftJoiningLetter(letter rune) bool {
	nonLeftJoiningLetters := [8]rune{'\u0627', '\u0622', '\u0698', '\u062F', '\u0630', '\u0631', '\u0632', '\u0648'}
	for _, item := range nonLeftJoiningLetters {
		if item == letter {
			return true
		}
	}
	return false
}

// IsPersianLetter checks if the letter is persian
func IsPersianLetter(ch rune) bool {
	return ch >= 0x600 && ch <= 0x6FF
}

// IsPersian checks if the input string contains persian Unicode only
func IsPersian(input string) bool {
	var isPersian = true
	for _, v := range input {
		if !unicode.IsSpace(v) && !IsPersianLetter(v) {
			isPersian = false
		}
	}
	return isPersian
}

// ToPersianDigits will convert english numbers to persian numbers in text
func ToPersianDigits(input string) string {
	return strings.NewReplacer(
		"0", "۰",
		"1", "۱",
		"2", "۲",
		"3", "۳",
		"4", "۴",
		"5", "۵",
		"6", "۶",
		"7", "۷",
		"8", "۸",
		"9", "۹",
	).Replace(input)
}

// ToEnglishDigits will convert persian numbers to english numbers in text
func ToEnglishDigits(input string) string {
	return strings.NewReplacer(
		"۰", "0",
		"۱", "1",
		"۲", "2",
		"۳", "3",
		"۴", "4",
		"۵", "5",
		"۶", "6",
		"۷", "7",
		"۸", "8",
		"۹", "9",
	).Replace(input)
}
