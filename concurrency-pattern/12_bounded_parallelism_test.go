package concurrency_pattern

import (
	"crypto/md5"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"
)

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func digester(done <-chan struct{}, paths <-chan string, c chan<- result) {
	for path := range paths {
		data, err := os.ReadFile(path)
		select {
		case c <- result{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}

func MD5ALL(root string) (map[string][md5.Size]byte, error) {
	done := make(chan struct{})
	defer close(done)
	paths, errc := walkFiles(done, root)
	c := make(chan result)
	var wg sync.WaitGroup
	const numDigesters = 20
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(done, paths, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			m[r.path] = r.sum
		}
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

func Test12(t *testing.T) {
	m, err := MD5ALL(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		log.Println(m[path], path)
	}
}
