package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func SendHTTPRequest(httpClient http.Client, url string, reqbody any, resbody any) error {
	method := http.MethodPost
	if reqbody == nil {
		method = http.MethodGet
	}

	encodedBody, err := json.Marshal(reqbody)
	if err != nil {
		return errors.Wrap(err, "couldn't json marshal request")
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(encodedBody))
	if err != nil {
		return errors.Wrap(err, "couldn't create request")
	}
	req.Header.Set("Content-Type", "application/json")

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

	}

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("status is %d instead of 200 content is %s", res.StatusCode, string(resdata)))
	}

	if resbody != nil {
		if err = json.Unmarshal(resdata, resbody); err != nil {
			return errors.Wrap(err, "couln't unmarshal data from json")
		}
	}

	return nil

}
