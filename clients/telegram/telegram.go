package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"tg-bot/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newToken(token),
		client:   http.Client{},
	}
}

func newToken(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest("getUpdates", q)
	if err != nil {
		return nil, e.Wrap("can't do updates", err)
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap("can't do updates", err)
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest("sendMessage", q)
	if err != nil {
		return e.Wrap("can't do sendMessage", err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {

	res := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, res.String(), nil)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}
	return body, nil
}
