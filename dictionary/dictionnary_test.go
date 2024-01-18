package dictionary

import (
	"reflect"
	"testing"
)

func TestAddEntry(t *testing.T) {

	d := New()
	word := "test_word"
	definition := "test_definition"
	d.Add(word, definition)

	entry, err := d.Get(word)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if entry.Definition != "test_definition" {
		t.Errorf("Expected definition 'test_definition', got %s", entry.Definition)
	}
}

func TestGetEntry(t *testing.T) {
	d := New()
	word := "test_word"
	definition := "test_definition"
	d.Add(word, definition)

	entry, err := d.Get("test_word")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if entry.Definition != "test_definition" {
		t.Errorf("Expected definition 'test_definition', got %s", entry.Definition)
	}
}

func TestRemoveEntry(t *testing.T) {
	d := New()
	word := "test_word"
	definition := "test_definition"
	d.Add(word, definition)

	// Supprimer l'entrée
	err := d.Remove(word)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Vérifier que l'entrée existe avant la suppression
	_, err = d.Get(word)
	if err == nil {
		t.Errorf("Expected error, entry should be removed")
	}
}

func TestListEntries(t *testing.T) {
	d := New()
	word1 := "apple"
	word2 := "banana"
	word3 := "cherry"
	definition1 := "a fruit"
	definition2 := "another fruit"
	definition3 := "yet another fruit"

	d.Add(word1, definition1)
	d.Add(word2, definition2)
	d.Add(word3, definition3)

	expectedWordList := []string{word1, word2, word3}
	expectedEntries := map[string]Entry{
		word1: {Word: word1, Definition: definition1},
		word2: {Word: word2, Definition: definition2},
		word3: {Word: word3, Definition: definition3},
	}

	wordList, entries := d.List()

	// Vérifier que la liste des mots est correcte
	if !reflect.DeepEqual(wordList, expectedWordList) {
		t.Errorf("Expected word list %v; got %v", expectedWordList, wordList)
	}

	// Vérifier que la liste des entrées est correcte
	if !reflect.DeepEqual(entries, expectedEntries) {
		t.Errorf("Expected entries %v; got %v", expectedEntries, entries)
	}
}
