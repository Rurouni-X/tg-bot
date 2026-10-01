package files

import (
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"tg-bot/lib/e"
	"tg-bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

func NewFileStorage(path string) Storage {
	return Storage{basePath: path}
}

func (s Storage) Save(page *storage.Page) error {
	path := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(path, 0774); err != nil {
		return e.Wrap("can't save page", err)
	}

	fname, err := fileName(page)
	if err != nil {
		return e.Wrap("can't save page", err)
	}

	path = filepath.Join(path, fname)

	file, err := os.Create(path)
	if err != nil {
		return e.Wrap("can't save page", err)
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap("can't save page", err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	fPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, e.Wrap("can't pick random page", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavePage
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return e.Wrap("can't remove file", err)
	}
	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {

	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't exist file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(path); {

	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap("can't exist file", err)

	}

	return true, nil
}

func (s Storage) decodePage(file string) (*storage.Page, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, e.Wrap("can't decode page file", err)
	}
	defer f.Close()

	var pg storage.Page

	if err := gob.NewDecoder(f).Decode(&pg); err != nil {
		return nil, e.Wrap("can't decode page file", err)
	}
	return &pg, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
