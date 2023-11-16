package dictionary

type Entry struct {
}

func (e Entry) String() string {

	return ""
}

type Dictionary struct {
	entries map[string]Entry
}

func New() *Dictionary {

	return nil
}

func (d *Dictionary) Add(word string, definition string) {

}

func (d *Dictionary) Get(word string) (Entry, error) {

	return Entry{}, nil
}

func (d *Dictionary) Remove(word string) {

}

func (d *Dictionary) List() ([]string, map[string]Entry) {

	return []string{}, d.entries
}
