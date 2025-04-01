package stdc

import (
	"strings"
	"unicode/utf8"
)

func CharsCount(str string) int {
	return utf8.RuneCountInString(str)
}

func CharCount(str string, char rune) int {
	count := 0

	for _, r := range str {
		if r == char {
			count += 1
		}
	}

	return count
}

func ReverseString(str *string) {
	runes := []rune(*str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	*str = string(runes)
}

func AddChar(str *string, char string, to_end bool) {
	if to_end {
		*str += char
	} else {
		*str = char + *str
	}
}

func RemoveChar(str *string, char rune) {
	var result strings.Builder

	for _, r := range *str {
		if r != char {
			result.WriteRune(r)
		}
	}

	*str = result.String()
}

func AddString(str *string, to_add string, to_end bool) {
	if to_end {
		*str += to_add
	} else {
		*str = to_add + *str
	}
}

func RemoveString(str *string, to_remove string, to_end bool) {
	if to_end {
		if strings.HasSuffix(*str, to_remove) {
			*str = (*str)[:len(*str)-len(to_remove)]
		}
	} else {
		*str = strings.ReplaceAll(*str, to_remove, "")
	}
}

func Replace(str *string, to_replace string, replaced string, replace_count byte) {
	*str = strings.Replace(*str, to_replace, replaced, int(replace_count))
}

func ReplaceAll(str *string, to_replace string, replaced string) {
	*str = strings.ReplaceAll(*str, to_replace, replaced)
}

func StringToArray(arr *[]string, str string) {
	*arr = make([]string, len(str))

	for i, r := range str {
		(*arr)[i] = string(r)
	}
}

func ArrayToString(str *string, arr []string) {
	for _, val := range arr {
		*str += val
	}
}
