package http2ftp

import (
	"net/http"
	"io/ioutil"
	"io"
	"bytes"
	"mime/multipart"
	"fmt"
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

func uploadFileData(url string, fieldName string, fileName string, reader io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	multipart, err := writer.CreateFormFile(fieldName, fileName)

	if err != nil {
		return err
	}

	io.Copy(multipart, reader)
	writer.Close()

	fmt.Println(">>>> Request")
	fmt.Println(body.String())

	request, err := http.NewRequest("POST", url, body)
	if (err != nil) {
		return err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	response, err := client.Do(request)

	if (err != nil) {
		return err
	}

	responseBody, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	fmt.Println(">>>> Response")
	fmt.Println(string(responseBody))

	return nil
}
