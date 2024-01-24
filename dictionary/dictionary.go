package dictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

const redisKey = "dictionary"

type Entry struct {
	Word       string
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

// RedisDictionary représente un dictionnaire stocké dans Redis.
type RedisDictionary struct {
	rdb   *redis.Client
	mutex sync.Mutex
}

type Dictionary struct {
	rdb *redis.Client // Redis client
}

func New(redisAddr string) *Dictionary {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})

	return &Dictionary{
		rdb: rdb,
	}
}

func NewRedisDictionary(redisAddr string, rdb *redis.Client) *RedisDictionary {
	// Assurez-vous que Redis est accessible et prêt à être utilisé.
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
	}

	return &RedisDictionary{
		rdb: rdb,
	}
}

// Add ajoute un mot et sa définition au dictionnaire Redis.
func (rd *RedisDictionary) Add(word, definition string) {
	entry := Entry{Word: word, Definition: definition}
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Error encoding entry:", err)
		return
	}

	rd.mutex.Lock()
	defer rd.mutex.Unlock()

	rd.rdb.HSet(context.Background(), "entries", word, data)
}

// Get récupère la définition d'un mot à partir du dictionnaire Redis.
func (rd *RedisDictionary) Get(word string) (Entry, error) {
	rd.mutex.Lock()
	defer rd.mutex.Unlock()

	data, err := rd.rdb.HGet(context.Background(), "entries", word).Result()
	if err == redis.Nil {
		return Entry{}, fmt.Errorf("word '%s' not found", word)
	} else if err != nil {
		fmt.Println("Error getting entry from Redis:", err)
		return Entry{}, err
	}

	var entry Entry
	err = json.Unmarshal([]byte(data), &entry)
	if err != nil {
		fmt.Println("Error decoding entry:", err)
		return Entry{}, err
	}

	return entry, nil
}

// List récupère la liste de tous les mots et de leurs définitions dans le dictionnaire Redis.
func (rd *RedisDictionary) List() ([]string, map[string]Entry) {
	rd.mutex.Lock()
	defer rd.mutex.Unlock()

	entries, err := rd.rdb.HGetAll(context.Background(), "entries").Result()
	if err != nil {
		fmt.Println("Error getting entries from Redis:", err)
		return nil, nil
	}

	wordList := make([]string, 0, len(entries))
	entryMap := make(map[string]Entry)

	for word, data := range entries {
		var entry Entry
		err := json.Unmarshal([]byte(data), &entry)
		if err != nil {
			fmt.Println("Error decoding entry:", err)
			continue
		}

		wordList = append(wordList, word)
		entryMap[word] = entry
	}

	return wordList, entryMap
}

// Remove supprime un mot du dictionnaire Redis.
func (rd *RedisDictionary) Remove(word string) error {
	rd.mutex.Lock()
	defer rd.mutex.Unlock()

	// Vérifier si le mot existe avant de le supprimer
	exists, err := rd.rdb.HExists(context.Background(), "entries", word).Result()
	if err != nil {
		fmt.Println("Error checking if word exists in Redis:", err)
		return err
	}

	if !exists {
		return fmt.Errorf("Word not found: %s", word)
	}

	// Supprimer le mot du dictionnaire
	rd.rdb.HDel(context.Background(), "entries", word)
	return nil
}
