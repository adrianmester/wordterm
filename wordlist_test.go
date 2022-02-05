package main

import (
	"testing"
	"time"
)

func TestWordList_GetWord(t *testing.T) {
	wl, err := NewWordList()
	if err != nil {
		t.Error(err)
	}
	word := wl.GetWord(227)
	if word != "those" {
		t.Error("Expected 'those' but got ", word)
	}
}

func TestWordList_GetWordByDate(t *testing.T) {
	wl, err := NewWordList()
	if err != nil {
		t.Error(err)
	}
	n, word := wl.GetWordByDate(time.Date(2022, 1, 31, 16, 0, 0, 0, time.Local))
	if n != 226 {
		t.Error("Expected 226 but got ", n)
	}
	if word != "light" {
		t.Error("Expected 'light' but got ", word)
	}
}
