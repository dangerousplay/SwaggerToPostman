package utils

import (
	"SwaggerToPostman/models"
	"encoding/json"
	"github.com/valyala/fastjson"
	"strings"
)

func ConvertSwaggerToPostman(bytes []byte) ([]byte, error) {

	parser := fastjson.Parser{}

	result, _ := parser.ParseBytes(bytes)

	info := result.GetObject("info")

	postman := models.Postman{
		Info: models.Info{
			Name:        info.Get("title").String(),
			Description: info.Get("description").String(),
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
	}

	items := []models.Item{
		{Name: "api"},
	}

	itemsIn := []models.ItemIn{}

	result.GetObject("paths").Visit(func(k []byte, v *fastjson.Value) {
		endpoint := string(k)
		for _, v := range convertRequest(v, endpoint, result) {
			itemsIn = append(itemsIn, v)
		}

	})

	items[0].Item = itemsIn

	postman.Item = items

	return json.Marshal(postman)
}

func convertRequest(value *fastjson.Value, endpoint string, all *fastjson.Value) []models.ItemIn {
	item := []models.ItemIn{}

	index := 0

	value.GetObject().Visit(func(k []byte, v *fastjson.Value) {
		itemin := models.ItemIn{}
		method := string(k)

		itemin.Request.Method = strings.ToUpper(method)
		itemin.Request.Description = removeScape(v.Get("description").String())
		itemin.Name = removeScape(v.Get("summary").String())
		itemin.Request.URL.Raw = "{{hostname}}" + endpoint

		path := &itemin.Request.URL.Path

		*path = []string{}

		for _, v := range strings.Split(endpoint, "/") {
			if v == "" {
				continue
			}

			*path = append(*path, v)
		}

		if is(method, GET) {
			h, q := convertParameters(value.Get("get").Get("parameters"))
			itemin.Request.Header = h

			itemin.Request.URL.Raw += q
		} else if is(method, POST) {

			rest := v.Get("requestBody").Get("content").Get("application/json").Get("schema").Get("$ref").String()

			bodu := recursiveGet(all, getPath(rest))

			ret := convertBodyPost(bodu)

			itemin.Request.Body = ret
		}

		item = append(item, itemin)
		index++
	})

	return item
}

func getPath(path string) []string {
	path = strings.Replace(path, "#/", "", 1)
	path = strings.Replace(path, `"`, "", -1)
	return strings.Split(strings.TrimSpace(path), "/")
}

func recursiveGet(value *fastjson.Value, path []string) *fastjson.Value {

	var val = value

	for _, v := range path {
		val = val.Get(v)
	}

	return val
}

func convertParameters(value *fastjson.Value) ([]models.Header, string) {
	headers := []models.Header{}
	parameters := []string{}

	for _, v := range value.GetArray() {
		in := removeScape(v.Get("in").String())
		name := removeScape(v.Get("name").String())

		if in == "header" {
			headers = append(headers, models.Header{
				Key: name,
			})
		} else if in == "query" {
			parameters = append(parameters, name)
		}
	}

	toReturn := ""

	if len(parameters) > 0 {
		toReturn = "?"
		for k, v := range parameters {
			if k == 0 {
				toReturn += v + "="
			} else {
				toReturn += "&" + v + "="
			}
		}
	}

	return headers, toReturn
}

func convertToInterfaceBody(value *fastjson.Value) map[string]interface{} {
	mapt := make(map[string]interface{})

	value.Get("properties").GetObject().Visit(func(key []byte, v *fastjson.Value) {
		exampleP := v.Get("example")
		var example string

		if exampleP != nil {
			example = removeScape(exampleP.String())
		}

		aType := v.Get("type")

		var typ = ""

		if aType != nil {
			typ = removeScape(aType.String())
		}

		if typ == "object" {
			mapt[string(key)] = convertToInterfaceBody(v)
		} else if typ == "array" {
			mapt[string(key)] = []string{example}
		} else {
			mapt[string(key)] = example
		}

	})

	return mapt
}

func convertBodyPost(value *fastjson.Value) models.Body {
	mapt := convertToInterfaceBody(value)

	body := models.Body{}

	body.Mode = "raw"

	bytes, err := json.Marshal(mapt)

	if err != nil {
		panic(err)
	}

	body.Raw = string(bytes)

	return body
}

func is(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

func removeScape(stirng string) string {
	return strings.Replace(stirng, `"`, "", -1)
}

const (
	POST   = "post"
	GET    = "get"
	DELETE = "delete"
	PUT    = "put"
)
