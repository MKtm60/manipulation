package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"manipulation/dictionary"

	"github.com/gorilla/mux"
)

const dictionaryFile = "dictionary.json"

type SaveData struct {
	Entries map[string]struct {
		Definition string `json:"definition"`
	} `json:"entries"`
}

func main() {
	router := mux.NewRouter()

	d := loadDictionary()

	router.HandleFunc("/add", func(w http.ResponseWriter, router *http.Request) {
		actionAdd(d, w, router)
	}).Methods("POST")

	router.HandleFunc("/define/{word}", func(w http.ResponseWriter, router *http.Request) {
		actionDefine(d, w, router)
	}).Methods("PUT")

	router.HandleFunc("/list", func(w http.ResponseWriter, router *http.Request) {
		actionList(d, w)
	}).Methods("GET")

	router.HandleFunc("/remove/{word}", func(w http.ResponseWriter, router *http.Request) {
		actionRemove(d, w, router)
	}).Methods("DELETE")

	router.HandleFunc("/exit", func(w http.ResponseWriter, router *http.Request) {
		saveDictionary(d)
		fmt.Fprintln(w, "Exit program.")
	}).Methods("GET")

	http.Handle("/", router)

	fmt.Println("Server started on :8090")
	http.ListenAndServe(":8090", router)
}

// Rest of the code remains unchanged...

func actionAdd(d *dictionary.Dictionary, w http.ResponseWriter, router *http.Request) {
	decoder := json.NewDecoder(router.Body)
	var requestData struct {
		Word       string `json:"word"`
		Definition string `json:"definition"`
	}

	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	d.Add(requestData.Word, requestData.Definition)
	response := fmt.Sprintf("Word '%s' added with definition '%s'.\n", requestData.Word, requestData.Definition)
	fmt.Fprintln(w, response)
}

func actionDefine(d *dictionary.Dictionary, w http.ResponseWriter, router *http.Request) {
	vars := mux.Vars(router)
	word := vars["word"]

	decoder := json.NewDecoder(router.Body)
	var requestData struct {
		Definition string `json:"definition"`
	}

	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	entry, err := d.Get(word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	d.Remove(word)
	d.Add(word, requestData.Definition)
	response := fmt.Sprintf("Definition of '%s': '%s'\n", entry.Word, entry.Definition)
	fmt.Fprintln(w, response)
}

func actionList(d *dictionary.Dictionary, w http.ResponseWriter) {
	words, entries := d.List()
	response := "Words in the dictionary:\n"
	for _, word := range words {
		response += fmt.Sprintf("%s: %s\n", word, entries[word])
	}
	fmt.Fprintln(w, response)
}

func actionRemove(d *dictionary.Dictionary, w http.ResponseWriter, router *http.Request) {
	vars := mux.Vars(router)
	word := vars["word"]

	entry, err := d.Get(word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	d.Remove(word)
	response := fmt.Sprintf("Word '%s' removed.\n", entry.Word)
	fmt.Fprintln(w, response)
}

func saveDictionary(d *dictionary.Dictionary) {
	saveData := SaveData{
		Entries: make(map[string]struct {
			Definition string `json:"definition"`
		}),
	}

	words, entries := d.List()
	for _, word := range words {
		saveData.Entries[word] = struct {
			Definition string `json:"definition"`
		}{Definition: entries[word].Definition}
	}

	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		fmt.Println("Error encoding dictionary:", err)
		return
	}

	err = os.WriteFile(dictionaryFile, data, 0644)
	if err != nil {
		fmt.Println("Error saving dictionary to file:", err)
	}
}

func loadDictionary() *dictionary.Dictionary {
	d := dictionary.New()

	data, err := os.ReadFile(dictionaryFile)
	if err != nil {
		fmt.Println("Error reading dictionary file:", err)
		return d
	}

	if len(data) == 0 {
		return d
	}

	var saveData SaveData
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		fmt.Println("Error decoding dictionary:", err)
		return d
	}

	for word, entry := range saveData.Entries {
		d.Add(word, entry.Definition)
	}

	return d
}
