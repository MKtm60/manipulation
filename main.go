package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"manipulation/dictionary"
	"manipulation/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

const (
	dictionaryFile = "dictionary.json"
	redisAddr      = "localhost:6379" // Adresse de votre serveur Redis
)

type SaveData struct {
	Entries map[string]struct {
		Definition string `json:"definition"`
	} `json:"entries"`
}

func main() {
	middleware.InitLogFile()

	// Initialisation du client Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})
	defer rdb.Close()

	router := mux.NewRouter()

	// Utilisation des middlewares
	router.Use(middleware.AuthenticationMiddleware)
	router.Use(middleware.LoggingToFileMiddleware)
	router.Use(middleware.ValidateDataMiddleware)

	// Création de l'instance du dictionnaire avec Redis
	d := dictionary.NewRedisDictionary(redisAddr, rdb)

	// Définition des routes
	router.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		actionAdd(d, w, r)
	}).Methods("POST")

	router.HandleFunc("/define/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionDefine(d, w, r)
	}).Methods("PUT")

	router.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		actionList(d, w)
	}).Methods("GET")

	router.HandleFunc("/remove/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionRemove(d, w, r)
	}).Methods("DELETE")

	router.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		saveDictionary(d)
		fmt.Fprintln(w, "Exit program.")
	}).Methods("GET")

	http.Handle("/", router)

	fmt.Println("Server started on :8090")
	http.ListenAndServe(":8090", router)
}

// Reste du code inchangé...

// Adaptation des fonctions pour utiliser RedisDictionary
func actionAdd(d *dictionary.RedisDictionary, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
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

func actionDefine(d *dictionary.RedisDictionary, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	decoder := json.NewDecoder(r.Body)
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

func actionList(d *dictionary.RedisDictionary, w http.ResponseWriter) {
	words, entries := d.List()
	response := "Words in the dictionary:\n"
	for _, word := range words {
		response += fmt.Sprintf("%s: %s\n", word, entries[word])
	}
	fmt.Fprintln(w, response)
}

func actionRemove(d *dictionary.RedisDictionary, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
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

func saveDictionary(d *dictionary.RedisDictionary) {
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
