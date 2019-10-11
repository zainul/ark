package ark

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// HTTP is caller
type HTTP struct {
	client *http.Client
}

// HTTPResult response struct
type HTTPResult struct {
	Response   []byte
	HTTPStatus int
}

// NewHTTP is http lib
func NewHTTP(timeout time.Duration) *HTTP {
	cl := http.Client{
		Timeout: timeout,
	}
	return &HTTP{
		client: &cl,
	}
}

// Call ..
func (c *HTTP) Call(
	method, url string,
	body io.Reader,
	header map[string]string,
) ([]byte, error) {
	req, err := c.GetRequest(method, url, body)

	if err != nil {
		log.Println("handle call error @http.Call", err)
		return nil, err
	}

	for key, val := range header {
		req.Header.Set(key, val)
	}

	resp, err := c.client.Do(req)

	if err != nil {
		log.Println("failed client do @http.Call")
		return nil, err
	}

	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("failed to read ioutil @http.Call", err)
		return nil, err
	}

	return resBody, nil
}

// GetRequest ...
func (c *HTTP) GetRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, body)

	if err != nil {
		log.Println("func GetRequest http @http.GetRequest", err)
		return nil, err
	}

	return req, nil
}
