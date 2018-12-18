package utils

import (
	"SwaggerToPostman/models"
	"encoding/json"
	"fmt"
	"github.com/valyala/fastjson"
	"strings"
)

func ConvertSwaggerToPostman(bytes []byte) ([]byte, error) {

	parser := fastjson.Parser{}

	result, _ := parser.ParseBytes(bytes)

	info := result.GetObject("info")

	postman := models.Postman {
		Info: models.Info{
			Name: info.Get("title").String(),
			Description: info.Get("description").String(),
			Schema: "https://schema.getpostman.com/collection/v2",
		},
	}

	items := []models.Item{}

	result.GetObject("paths").Visit(func(k []byte, v *fastjson.Value) {
		items = append(items,convertRequest(v, result))
	})

	return json.Marshal(postman)
}

func convertRequest( value *fastjson.Value, all *fastjson.Value) models.Item {
	item := models.Item{}

	itemin := models.ItemIn{}
	value.GetObject().Visit(func(k []byte, v *fastjson.Value) {
		method := string(k)

		itemin.Request.Method = method
		itemin.Request.Description = v.Get("description").String()
		itemin.Name = v.Get("summary").String()

		if is(method, GET) {

		} else if is(method, POST) {

			rest := v.Get("requestBody").Get("content").Get("application/json").Get("schema").Get("$ref").String()

			bodu := recursiveGet(all, getPath(rest))

			ret := convertBodyPost(bodu)
			fmt.Println(ret)
		}
	})

	return item
}

func getPath(path string) []string {
	path = strings.Replace(path, "#/", "", 1)
	path = strings.Replace(path, `"`,"",-1)
	return strings.Split(strings.TrimSpace(path), "/")
}

func recursiveGet(value *fastjson.Value,path []string) *fastjson.Value {

	var val = value

	for _, v := range path {
		val = val.Get(v)
	}

	return val
}

func convertBodyPost(value *fastjson.Value) models.Body {
	mapt := make(map[string]string)

	body := models.Body{}

	value.Get("properties").GetObject().Visit(func(key []byte, v *fastjson.Value) {
		example := v.Get("example").String()
		example = strings.Replace(example, `"`, "", -1)

		mapt[string(key)] = example
	})

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

const (
	POST = "post"
	GET = "get"
	DELETE = "delete"
	PUT = "put"
)