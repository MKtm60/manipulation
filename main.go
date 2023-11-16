package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"manipulation/dictionary"
)

func main() {
	dic := dictionary.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an action:")
		fmt.Println("1. Add")
		fmt.Println("2. Define")
		fmt.Println("3. Remove")
		fmt.Println("4. List")
		fmt.Println("5. Exit")

		choice, err := getUserChoice(reader)
		if err != nil {
			fmt.Println("Choice Error:", err)
			continue
		}

		switch choice {
		case 1:
			actionAdd(dic, reader)
		case 2:
			actionDefine(dic, reader)
		case 3:
			actionRemove(dic, reader)
		case 4:
			actionList(dic)
		case 5:
			fmt.Println("Exit program.")
			return
		default:
			fmt.Println("Invalid choice. Please choose an option valid.")
		}
	}
}

func getUserChoice(reader *bufio.Reader) (int, error) {
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

func actionAdd(dic *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Print("Enter the definition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	dic.Add(word, definition)
	fmt.Printf("Word '%s' added with definition '%s'.\n", word, definition)
}

func actionDefine(dic *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := dic.Get(word)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Definition of '%s': %s\n", entry.Word, entry.Definition)
}

func actionRemove(dic *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	dic.Remove(word)
	fmt.Printf("Word '%s' removed.\n", word)
}

func actionList(dic *dictionary.Dictionary) {
	words, entries := dic.List()
	fmt.Println("Words in the dictionary:")
	for _, word := range words {
		fmt.Printf("%s: %s\n", word, entries[word])
	}
}
