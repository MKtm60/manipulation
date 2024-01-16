package main

import (
	"bufio"
	"fmt"
	"os"

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

		choice, err := dictionary.GetUserChoice(reader)
		if err != nil {
			fmt.Println("Choice Error:", err)
			continue
		}

		switch choice {
		case 1:
			dictionary.ActionAdd(dic, reader)
		case 2:
			dictionary.ActionDefine(dic, reader)
		case 3:
			dictionary.ActionRemove(dic, reader)
		case 4:
			dictionary.ActionList(dic)
		case 5:
			fmt.Println("Exit program.")
			return
		default:
			fmt.Println("Invalid choice. Please choose an option valid.")
		}
	}
}
