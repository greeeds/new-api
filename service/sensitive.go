package service

import (
	"bytes"
	"fmt"
	"github.com/anknown/ahocorasick"
	"one-api/constant"
	"strings"
)

const DEFAULT_WORD = "**###**"
const SPLIT_KEY = "$$"

// SensitiveWordContains 是否包含敏感词，返回是否包含敏感词和敏感词列表
func SensitiveWordContains(text string) (bool, []string) {
	if len(constant.SensitiveWords) == 0 {
		return false, nil
	}
	checkText := strings.ToLower(text)
	// 构建一个AC自动机
	m := initAc()
	hits := m.MultiPatternSearch([]rune(checkText), false)
	if len(hits) > 0 {
		words := make([]string, 0)
		for _, hit := range hits {
			words = append(words, string(hit.Word))
		}
		return true, words
	}
	return false, nil
}

// SensitiveWordReplace 敏感词替换，返回是否包含敏感词和替换后的文本
func SensitiveWordReplace(text string, returnImmediately bool) (bool, []string, string) {
	if len(constant.SensitiveWords) == 0 {
		return false, nil, text
	}
	checkText := strings.ToLower(text)
	m := initAc()
	hits := m.MultiPatternSearch([]rune(checkText), returnImmediately)
	replaceWordMap := replaceMap()
	if len(hits) > 0 {
		textRunes := []rune(text)
		words := make([]string, 0, len(hits))
		var builder strings.Builder
		posOffset := 0
		for _, hit := range hits {
			pos := hit.Pos
			word := string(hit.Word)
			replaceWord := replaceWordMap[word]
			builder.WriteString(string(textRunes[posOffset:pos]))
			builder.WriteString(replaceWord)
			posOffset = pos + len([]rune(word))
			words = append(words, word)
		}
		builder.WriteString(string(textRunes[posOffset:]))
		return true, words, builder.String()
	}
	return false, nil, text
}

func initAc() *goahocorasick.Machine {
	m := new(goahocorasick.Machine)
	dict := readRunes()
	if err := m.Build(dict); err != nil {
		fmt.Println(err)
		return nil
	}
	return m
}

func readRunes() [][]rune {
	var dict [][]rune

	for _, word := range constant.SensitiveWords {
		word = strings.ToLower(strings.Split(word, SPLIT_KEY)[0])
		l := bytes.TrimSpace([]byte(word))
		dict = append(dict, bytes.Runes(l))
	}

	return dict
}

func replaceMap() map[string]string {
	result := make(map[string]string)
	for _, word := range constant.SensitiveWords {
		if strings.Contains(word, SPLIT_KEY) {
			parts := strings.Split(word, SPLIT_KEY)
			if len(parts) == 2 {
				result[strings.ToLower(parts[0])] = parts[1]
			} else {
				result[strings.ToLower(parts[0])] = ""
			}
		} else {
			result[strings.ToLower(word)] = DEFAULT_WORD
		}
	}
	return result
}
