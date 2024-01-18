package dictionary

import (
	"fmt"
	"sync"
)

type Entry struct {
	Word       string
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	entries       map[string]Entry
	updateChannel chan dictionaryUpdate
	mutex         sync.Mutex // guards
}

func New() *Dictionary {
	d := &Dictionary{
		entries:       make(map[string]Entry),
		updateChannel: make(chan dictionaryUpdate),
	}

	go d.startConcurrentOperations()

	return d
}

type dictionaryUpdate struct {
	entry Entry
	del   bool
}

func (d *Dictionary) startConcurrentOperations() {
	for {
		update := <-d.updateChannel

		d.mutex.Lock()

		if update.del {
			delete(d.entries, update.entry.Word)
		} else {
			d.entries[update.entry.Word] = update.entry
		}

		d.mutex.Unlock()
	}
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Word: word, Definition: definition}
	d.entries[word] = entry
	d.updateChannel <- dictionaryUpdate{entry: entry}
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, found := d.entries[word]
	d.updateChannel <- dictionaryUpdate{entry: entry}
	if !found {
		return Entry{}, fmt.Errorf("word '%s' not found", word)
	}
	return entry, nil
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	wordList := make([]string, 0, len(d.entries))
	for word := range d.entries {
		wordList = append(wordList, word)
	}
	return wordList, d.entries
}

func (d *Dictionary) Remove(word string) error {
	// Signal de suppression via le canal
	d.updateChannel <- dictionaryUpdate{entry: Entry{Word: word}, del: true}

	_, found := d.entries[word]
	if !found {
		return fmt.Errorf("Word not found: %s", word)
	}

	// Effectuer la suppression dans le dictionnaire

	delete(d.entries, word)
	return nil
}
