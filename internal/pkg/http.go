package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type header struct {
	key string
	val string
}

func WithHeader(key, val string) header {
	return header{
		key: key,
		val: val,
	}
}

func SendHTTPRequest(httpClient http.Client, url string, reqbody any, resbody any, headers ...header) error {
	logger := GetLogger("http_request").With(zap.String("package", "pkg"), zap.String("module", "http"))

	method := http.MethodPost
	if reqbody == nil {
		method = http.MethodGet
	}

	encodedBody, err := json.Marshal(reqbody)
	if err != nil {
		return errors.Wrap(err, "couldn't json marshal request")
	}
	logger.Debug("request body: " + string(encodedBody))

	req, err := http.NewRequest(method, url, bytes.NewReader(encodedBody))
	if err != nil {
		return errors.Wrap(err, "couldn't create request")
	}
	req.Header.Set("Content-Type", "application/json")
	for _, header := range headers {
		req.Header.Set(header.key, header.val)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "couldn't send request")
	}

	var resdata []byte
	if res.Body != nil {
		defer res.Body.Close()
		resdata, err = io.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "couldn't read response")
		}
		logger.Debug("response body: " + string(resdata))
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("status is %d instead of 200 content is %s", res.StatusCode, string(resdata)))
	}

	switch v := resbody.(type) {
	case *string:
		*v = string(resdata)

	default:
		if resbody != nil {
			if err = json.Unmarshal(resdata, resbody); err != nil {
				return errors.Wrap(err, "couln't unmarshal data from json")
			}
		}
	}

	return nil

}
