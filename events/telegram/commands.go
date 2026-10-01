package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"tg-bot/lib/e"
	"tg-bot/storage"
)

const (
	RndCmd   = "/rnd"
	StartCmd = "/start"
	HelpCmd  = "/help"
)

func (p *Processor) doCmd(chatID int, text, userName string) error {

	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, userName)

	if isURL(text) {
		return p.savePage(chatID, text, userName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, userName)
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageUrl, userName string) error {
	page := storage.Page{
		URL:      pageUrl,
		UserName: userName,
	}
	isExist, err := p.storage.IsExist(&page)
	if err != nil {
		return e.Wrap("can't do command save page", err)
	}

	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(&page); err != nil {
		return e.Wrap("can't do command save page", err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return e.Wrap("can't do command save page", err)
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, userName string) error {

	page, err := p.storage.PickRandom(userName)
	if err != nil && errors.Is(err, storage.ErrNoSavePage) {
		return e.Wrap("can't do command send random", err)
	}

	if errors.Is(err, storage.ErrNoSavePage) {
		p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return e.Wrap("can't do command send random", err)
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
