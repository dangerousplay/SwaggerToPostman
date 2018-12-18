package main

import (
	"SwaggerToPostman/utils"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

func main() {

	file, err := os.Open("swagger.yaml")

	checkError(err)

	buf, err := ioutil.ReadAll(file)

	checkError(err)

	var u = make(map[string]interface{})

	yaml.Unmarshal(buf,&u)

	bufj, err := yaml.YAMLToJSON(buf)

	post, err := utils.ConvertSwaggerToPostman(bufj)

	checkError(err)

	err = ioutil.WriteFile("postman.json", post, os.FileMode(os.O_CREATE))

	checkError(err)
}

func checkError(err error){
	if err != nil {
		panic(err)
	}
}



