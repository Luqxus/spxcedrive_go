package main

import "log"

func main() {

	fsWatcher, err := NewWatcher("./dir")
	if err != nil {
		log.Fatal(err)
	}

	if err := fsWatcher.Watch(); err != nil {
		log.Fatal(err)
	}

}
