package logx

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"
)

type HttpWriter struct {
	url string
}

func NewHttpWriter(url string) *HttpWriter {
	return &HttpWriter{
		url: url,
	}
}

func (w *HttpWriter) Write(p []byte) (n int, err error) {
	if len(p) > 0 {
		return w.request(p)
	}
	return -1, nil
}

func (w *HttpWriter) request(p []byte) (int, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    100,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	p = p[:len(p)-1]
	body := bytes.NewBuffer(p)
	req, err := http.NewRequest("POST", w.url, body)
	if err != nil {
		return -1, err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	_, err = client.Do(req)
	if err != nil {
		return -1, err
	}
	return 0, nil
}
