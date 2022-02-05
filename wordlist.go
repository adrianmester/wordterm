package main

import (
	_ "embed"
	"encoding/json"
	"time"
)

type WordList struct {
	Targets  []string `json:"targets"`
	Valid    []string `json:"valid"`
	allWords map[string]struct{}
}

//go:embed words.json
var wordFile []byte

func NewWordList() (WordList, error) {
	var wl WordList
	err := json.Unmarshal(wordFile, &wl)
	if err != nil {
		return wl, err
	}
	wl.allWords = make(map[string]struct{})
	for _, word := range wl.Valid {
		wl.allWords[word] = struct{}{}
	}
	for _, word := range wl.Targets {
		wl.allWords[word] = struct{}{}
	}
	return wl, nil
}

func (wl WordList) CheckWord(word string) bool {
	_, ok := wl.allWords[word]
	return ok
}

func (wl WordList) GetWord(n int64) string {
	// roll over after the last word
	n = n % int64(len(wl.Targets))
	return wl.Targets[n]
}

func (wl WordList) GetWordByDate(date time.Time) (int64, string) {
	today := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	wordleEpoch := time.Date(2021, 06, 19, 0, 0, 0, 0, time.UTC)
	n := today.Sub(wordleEpoch).Hours() / 24
	return int64(n), wl.GetWord(int64(n))
}
