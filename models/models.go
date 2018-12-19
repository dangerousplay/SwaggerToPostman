package models

type Info struct {
	PostmanID   string `json:"_postman_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
}

type Url struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type Body struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Request struct {
	Method      string   `json:"method"`
	Header      []Header `json:"header"`
	Body        Body     `json:"body"`
	URL         Url      `json:"url"`
	Description string   `json:"description"`
}

type ItemIn struct {
	Name     string     `json:"name"`
	Request  Request    `json:"request"`
	Response []Response `json:"response"`
}

type Item struct {
	Name        string   `json:"name"`
	Item        []ItemIn `json:"item"`
	Description string   `json:"description"`
}

type Postman struct {
	Info Info   `json:"info"`
	Item []Item `json:"item"`
}

type Response struct {
	Body string
	Code int
}
