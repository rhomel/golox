package vm

// Since Go already has a built-in hash table (map) and hashtables aren't
// really the big focus of what I want to learn, I will just lean on the Go map
// for all hash-table functions. To stay in sync with the book however we will
// try to duplicate the functions so the code that utilizes the hash table is
// more or less similar.
//
// One key part of the hash-table chapter however is string interning (or
// string de-duplication). The table here will also do string interning since
// Go does not provide this for us with dynamically allocated strings.

type Entry struct {
	key   *ObjectString
	value Value
}

type Table struct {
	entries map[*ObjectString]*Entry
}

func (t *Table) initTable() {
	t.entries = make(map[*ObjectString]*Entry)
}

func (t *Table) freeTable() {
	t.initTable()
}

func (t *Table) Set(key *ObjectString, value Value) bool {
	entry, ok := t.entries[key]
	isNewKey := !ok
	if isNewKey {
		entry = &Entry{}
		t.entries[key] = entry
	}
	entry.key = key
	entry.value = value
	return isNewKey
}

func (t *Table) FindEntry(key *ObjectString) *Entry {
	for _, e := range t.entries {
		if e.key.String == key.String {
			return e
		}
	}
	return nil
}

func (t *Table) Get(key *ObjectString, value *Value) bool {
	entry := t.FindEntry(key)
	if entry == nil {
		return false
	}
	*value = entry.value
	return true
}

func (t *Table) Delete(key *ObjectString) bool {
	entry := t.FindEntry(key)
	if entry == nil {
		return false
	}
	delete(t.entries, entry.key)
	return true
}

func (from *Table) AddAll(to *Table) {
	for _, e := range from.entries {
		to.Set(e.key, e.value)
	}
}

func (t *Table) FindString(chars string, hash uint32) *ObjectString {
	for _, e := range t.entries {
		if len(e.key.String) == len(chars) && e.key.Hash == hash && e.key.String == chars {
			return e.key
		}
	}
	return nil
}
