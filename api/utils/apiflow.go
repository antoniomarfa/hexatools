package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type FlowApi struct {
	secretKey string
	apiURL    string
	apiKey    string
}

func (f *FlowApi) SetApiKey(apiKey string) {
	f.apiKey = apiKey
}

func (f *FlowApi) SetSecretKey(secretKey string) {
	f.secretKey = secretKey
}

func (f *FlowApi) SetApiURL(apiURL string) {
	f.apiURL = apiURL
}

func (f *FlowApi) Send(service string, params map[string]string, method string) (map[string]interface{}, error) {
	method = strings.ToUpper(method)
	url := f.apiURL + "/" + service
	params["apiKey"] = f.apiKey

	data := f.getPack(params, method)
	sign := f.sign(params)

	var response map[string]interface{}
	var err error

	if method == "GET" {
		response, err = f.httpGet(url, data, sign)
	} else {
		response, err = f.httpPost(url, data, sign)
	}

	if err != nil {
		return nil, err
	}

	if body, ok := response["body"].(map[string]interface{}); ok {
		if code, ok := response["code"].(int); ok {
			if code == 200 {
				return body, nil
			} else if code == 400 || code == 401 {
				return nil, errors.New(body["message"].(string))
			} else {
				return nil, fmt.Errorf("unexpected error occurred.http_code: %d", code)
			}
		}
	}

	return nil, errors.New("unexpected error occurred")
}

func (f *FlowApi) getPack(params map[string]string, method string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var data strings.Builder
	for _, key := range keys {
		if method == "GET" {
			data.WriteString("&" + url.QueryEscape(key) + "=" + url.QueryEscape(params[key]))
		} else {
			data.WriteString("&" + key + "=" + params[key])
		}
	}
	return data.String()[1:]
}

func (f *FlowApi) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var toSign strings.Builder
	for _, key := range keys {
		toSign.WriteString("&" + key + "=" + params[key])
	}
	toSignStr := toSign.String()[1:]

	h := hmac.New(sha256.New, []byte(f.secretKey))
	h.Write([]byte(toSignStr))
	return hex.EncodeToString(h.Sum(nil))
}

func (f *FlowApi) httpGet(url, data, sign string) (map[string]interface{}, error) {
	fullURL := url + "?" + data + "&s=" + sign
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"body": result,
		"code": resp.StatusCode,
	}, nil
}

func (f *FlowApi) httpPost(urlStr, data, sign string) (map[string]interface{}, error) {
	fullData := data + "&s=" + sign

	// Crear un mapa url.Values y asignar fullData a la clave "data"
	formData := url.Values{}
	formData.Set("data", fullData)

	resp, err := http.PostForm(urlStr, formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"body": result,
		"code": resp.StatusCode,
	}, nil
}
