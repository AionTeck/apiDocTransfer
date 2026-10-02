package yaml_doc_processor

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

func GetYamlDoc(client *http.Client, url, docName, authToken string) (processedDoc string, err error) {
	yamlDoc, err := getYamlDoc(client, url, authToken)
	if err != nil {
		return
	}

	processedDoc, err = changeYamlDocTitle(yamlDoc, docName)
	if err != nil {
		return
	}

	return
}

func getYamlDoc(client *http.Client, url, authToken string) (yamlDoc string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", authToken)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("Invalid status code: " + resp.Status)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	yamlDoc = string(bodyBytes)

	return
}

func changeYamlDocTitle(src, newTitle string) (result string, err error) {
	file, err := parser.ParseBytes([]byte(src), parser.ParseComments)
	if err != nil {
		return
	}

	p, err := yaml.PathString("$.info.title")
	if err != nil {
		return
	}

	val, err := yaml.Marshal(newTitle)
	if err != nil {
		return
	}

	if err = p.ReplaceWithReader(file, strings.NewReader(string(val))); err != nil {
		return
	}

	result = file.String()

	return
}
