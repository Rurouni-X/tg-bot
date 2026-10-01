package telegram

import (
	"errors"
	"tg-bot/clients/telegram"
	"tg-bot/events"
	"tg-bot/lib/e"
	"tg-bot/storage"
)

var (
	ErrUnknowEvent    = errors.New("unknow event type")
	ErrUnknowMetaType = errors.New("unknow meta type")
)

type Meta struct {
	ChatID   int
	UserName string
}

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func NewProcessor(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {

	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(update) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(update))

	for _, u := range update {
		res = append(res, event(u))
	}
	p.offset = update[len(update)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {

	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("cant't process message", ErrUnknowEvent)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(meta.ChatID, event.Text, meta.UserName); err != nil {
		return e.Wrap("can't process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknowMetaType)
	}
	return res, nil
}

func event(u telegram.Update) events.Event {

	updType := fetchType(u)

	res := events.Event{
		Type: updType,
		Text: fetchText(u),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			UserName: u.Message.From.Username,
		}
	}
	return res
}

func fetchType(u telegram.Update) events.Type {

	if u.Message == nil {
		return events.Unknow
	}
	return events.Message
}

func fetchText(u telegram.Update) string {

	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}
