package dictionary

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Entry struct {
	Word       string
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	entries map[string]Entry
}

func GetUserChoice(reader *bufio.Reader) (int, error) {
	fmt.Print("Enter your choice: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	text = strings.TrimSpace(text)
	choice, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	return choice, nil
}

func New() *Dictionary {
	return &Dictionary{
		entries: make(map[string]Entry),
	}
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Word: word, Definition: definition}
	d.entries[word] = entry
}

func ActionAdd(dic *Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Print("Enter the definition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	dic.Add(word, definition)
	fmt.Printf("Word '%s' added with definition '%s'.\n", word, definition)
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, found := d.entries[word]
	if !found {
		return Entry{}, fmt.Errorf("word '%s' not found", word)
	}
	return entry, nil
}

func ActionDefine(dic *Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := dic.Get(word)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Enter the new definition : ")
	newDefinition, _ := reader.ReadString('\n')
	newDefinition = strings.TrimSpace(newDefinition)
	dic.Remove(word)
	dic.Add(word, newDefinition)
	fmt.Printf("Definition of '%s' : '%s\n'", entry.Word, entry.Definition)
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func ActionRemove(dic *Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := dic.Get(word)
	if err != nil {
		fmt.Println(err)
		return
	}

	dic.Remove(word)
	fmt.Printf("Word '%s' removed.\n", entry.Word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	wordList := make([]string, 0, len(d.entries))
	for word := range d.entries {
		wordList = append(wordList, word)
	}
	return wordList, d.entries
}

func ActionList(dic *Dictionary) {
	words, entries := dic.List()
	fmt.Println("Words in the dictionary:")
	for _, word := range words {
		fmt.Printf("%s: %s\n", word, entries[word])
	}
}
