package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func makeKey() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func readRequestBody(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	if _, err = url.Parse(string(body)); err != nil {
		Logger.Error(err)
		return "", err
	}
	return string(body), nil
}

func readURLFromAPIRequestBody(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	var req ShortenReq
	if err = json.Unmarshal(body, &req); err != nil {
		return "", err
	}

	if _, err = url.Parse(req.URL); err != nil {
		Logger.Error(err)
		return "", err
	}
	return req.URL, nil
}
