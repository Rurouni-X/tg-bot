package storage

import (
	"crypto/sha1"
	"errors"
	"io"
	"tg-bot/lib/e"
)

var ErrNoSavePage = errors.New("no saved page")

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("hash page", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("hash page", err)
	}

	return string(h.Sum(nil)), nil
}
