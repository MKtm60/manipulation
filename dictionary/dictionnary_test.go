package dictionary

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisDictionary_Add(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rd := NewRedisDictionary("localhost:6379", rdb)

	word := "test_word"
	definition := "test_definition"

	rd.Add(word, definition)

	entry, err := rd.Get(word)
	assert.NoError(t, err)
	assert.Equal(t, word, entry.Word)
	assert.Equal(t, definition, entry.Definition)
}

func TestRedisDictionary_Get(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rd := NewRedisDictionary("localhost:6379", rdb)

	word := "test_word"
	definition := "test_definition"

	rd.Add(word, definition)

	entry, err := rd.Get(word)
	assert.NoError(t, err)
	assert.Equal(t, word, entry.Word)
	assert.Equal(t, definition, entry.Definition)

	// Test avec un mot qui n'existe pas
	nonexistentWord := "nonexistent_word"
	_, err = rd.Get(nonexistentWord)
	assert.Error(t, err)
	assert.EqualError(t, err, "word 'nonexistent_word' not found")
}

func TestRedisDictionary_List(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rd := NewRedisDictionary("localhost:6379", rdb)

	word1 := "test_word1"
	definition1 := "test_definition1"

	word2 := "test_word2"
	definition2 := "test_definition2"

	rd.Add(word1, definition1)
	rd.Add(word2, definition2)

	wordList, entryMap := rd.List()

	assert.Contains(t, wordList, word1)
	assert.Contains(t, wordList, word2)

	entry1, exists1 := entryMap[word1]
	entry2, exists2 := entryMap[word2]

	assert.True(t, exists1)
	assert.True(t, exists2)

	assert.Equal(t, definition1, entry1.Definition)
	assert.Equal(t, definition2, entry2.Definition)
}

func TestRedisDictionary_Remove(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rd := NewRedisDictionary("localhost:6379", rdb)

	word := "test_word"
	definition := "test_definition"

	rd.Add(word, definition)

	// Vérifier que le mot existe avant la suppression
	_, err := rd.Get(word)
	assert.NoError(t, err)

	err = rd.Remove(word)
	assert.NoError(t, err)

	// Vérifier que le mot a été supprimé
	_, err = rd.Get(word)
	assert.Error(t, err)
	assert.EqualError(t, err, "word 'test_word' not found")

	// Test avec un mot qui n'existe pas
	nonexistentWord := "nonexistent_word"
	err = rd.Remove(nonexistentWord)
	assert.Error(t, err)
	assert.EqualError(t, err, "Word not found: nonexistent_word")
}
