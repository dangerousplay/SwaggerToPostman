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
			Schema:      "https://schema.getpostman.com/collection/v2",
		},
	}

	items := []models.Item{}

	result.GetObject("paths").Visit(func(k []byte, v *fastjson.Value) {
		items = append(items, convertRequest(v, result))
	})

	return json.Marshal(postman)
}

func convertRequest(value *fastjson.Value, all *fastjson.Value) models.Item {
	item := models.Item{}

	value.GetObject().Visit(func(k []byte, v *fastjson.Value) {
		itemin := models.ItemIn{}
		method := string(k)

		itemin.Request.Method = method
		itemin.Request.Description = removeScape(v.Get("description").String())
		itemin.Name = removeScape(v.Get("summary").String())

		if is(method, GET) {
			h, q := convertParameters(value.Get("get").Get("parameters"))
			itemin.Request.Header = h
		} else if is(method, POST) {

			rest := v.Get("requestBody").Get("content").Get("application/json").Get("schema").Get("$ref").String()

			bodu := recursiveGet(all, getPath(rest))

			ret := convertBodyPost(bodu)

			itemin.Request.Body = ret
		}
		item.Item = append(item.Item, itemin)
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
	parameters := "?"

	for k, v := range value.GetArray() {
		in := removeScape(v.Get("in").String())
		name := removeScape(v.Get("name").String())

		if in == "header" {
			headers = append(headers, models.Header{
				Key: name,
			})
		} else if in == "query" {
			if k >= len(value.GetArray()) {
				parameters += name
			} else {
				parameters += name + "=&"
			}
		}
	}

	return headers, parameters
}

func convertToInterfaceBody(value *fastjson.Value) map[string]interface{} {
	mapt := make(map[string]interface{})

	value.Get("properties").GetObject().Visit(func(key []byte, v *fastjson.Value) {
		exampleP := v.Get("example")
		var example string

		if exampleP != nil {
			example = removeScape(exampleP.String())
		}

		typ := removeScape(v.Get("type").String())

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

	body.Mode = "application/json"

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
