package http2ftp

import (
	"net/http"
	"io/ioutil"
)

func fetchHTTPResourceBody(method string, url string) (string, error) {
	client := http.Client{}
	request, clientError := http.NewRequest(method, url, nil)

	if clientError != nil {
		return "", clientError
	}

	response, requestError := client.Do(request)

	if (requestError != nil) {
		return "", requestError
	}

	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	return string(body), nil
}
