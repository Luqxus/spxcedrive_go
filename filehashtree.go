package main

import "fmt"

type FileHashTree struct {
	tree map[string]string
}

func (t *FileHashTree) Add(k, v string) {
	value, ok := t.tree[k]
	if ok {
		if value == v {
			return
		}
	}

	t.tree[k] = v
}

func (t *FileHashTree) Remove(k string) error {
	_, ok := t.tree[k]
	if !ok {
		return fmt.Errorf("file (%s) not found", k)
	}

	delete(t.tree, k)

	return nil
}
