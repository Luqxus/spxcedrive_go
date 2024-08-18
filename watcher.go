package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type FsWatcher struct {
	dir     string
	watcher *fsnotify.Watcher
	hasher  Hasher
}

func NewWatcher(dir string) (*FsWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &FsWatcher{
		dir:     dir,
		watcher: w,
		hasher:  &DefaultHasher{},
	}, nil
}

func (fs *FsWatcher) Watch() error {
	defer fs.watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-fs.watcher.Events:
				if !ok {
					return
				}

				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}

				f, err := os.Open(event.Name)
				if err != nil {
					log.Println(err)
					continue
				}

				h, err := fs.hasher.Hash(f)
				if err != nil {
					log.Println(err)
					continue
				}

				fmt.Printf("File sum : %s", hex.EncodeToString(h[:]))

			case err, ok := <-fs.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := fs.watcher.Add(fs.dir)
	if err != nil {
		return err
	}

	<-make(chan struct{})

	return nil
}
