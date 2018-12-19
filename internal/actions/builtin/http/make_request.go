package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
)

type requestInput struct {
	URL     string				`json:"url"`
	Method  string				`json:"method"`
	Body    string				`json:"body"`
	Headers map[string]string	`json:"headers"`
}

type responseOutput struct {
	Status  int					`json:"status"`
	Headers map[string][]string	`json:"headers"`
	Body    string				`json:"body"`
}

func outFromHttpResponse(response *http.Response) (responseOutput, error) {
	bodyBytes, err := ioutil.ReadAll(response.Body)

	return responseOutput{
		Status:  response.StatusCode,
		Headers: response.Header,
		Body:    string(bodyBytes),
	}, err
}

func (ri requestInput) toGoRequest() (*http.Request, error) {
	r, err := http.NewRequest(ri.Method, ri.URL, strings.NewReader(ri.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range ri.Headers {
		r.Header.Add(k, v)
	}

	return r, nil
}

func makeRequest(ri requestInput, ctx context.Context) (responseOutput, error) {
	r, err := ri.toGoRequest()
	if err != nil {
		return responseOutput{}, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return responseOutput{}, err
	}

	return outFromHttpResponse(resp)
}
