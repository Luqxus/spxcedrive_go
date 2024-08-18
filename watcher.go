package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type FsWatcher struct {
	dir               string
	watcher           *fsnotify.Watcher
	hasher            Hasher
	PathTransformFunc PathTransformFunc
	HashTree          *FileHashTree
}

func NewWatcher(dir string) (*FsWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &FsWatcher{
		dir:               dir,
		watcher:           w,
		hasher:            &DefaultHasher{},
		PathTransformFunc: DefaultPathTransformFunc,
		HashTree:          &FileHashTree{tree: make(map[string]string)},
	}, nil
}

func (fs *FsWatcher) Watch() error {
	defer fs.watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-fs.watcher.Events:

				fmt.Println(event.String())
				if !ok {
					return
				}

				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}

				f, err := os.Open(event.Name)
				if err != nil {
					fs.deleteFile(event.Name)
					continue
				}

				if err := fs.updateFile(f, event.Name); err != nil {
					log.Println(err)
				}

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

func (fs *FsWatcher) deleteFile(path string) {
	key := fs.PathTransformFunc(path)
	if err := fs.HashTree.Remove(key); err != nil {
		fmt.Printf("failed to delete file with error : %s", err.Error())
		return
	}

	fmt.Printf("file (%s) deleted.\n", path)
}

func (fs *FsWatcher) updateFile(f io.Reader, path string) error {
	h, err := fs.hasher.Hash(f)
	if err != nil {
		return nil
	}

	// transform path file & add to file hash tree
	key := fs.PathTransformFunc(path)
	fs.HashTree.Add(key, hex.EncodeToString(h))
	fmt.Printf("File sum : %s", hex.EncodeToString(h[:]))

	fmt.Printf("Hash Tree : %+v\n", fs.HashTree.tree)

	return nil
}
