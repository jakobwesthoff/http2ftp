package http2ftp

import "net/http"

func doHTTPRequest(method string, url string) (*http.Response, error) {
	client := http.Client{}
	request, clientError := http.NewRequest(method, url, nil)

	if clientError != nil {
		return nil, clientError
	}

	response, requestError := client.Do(request)

	return response, requestError
}
