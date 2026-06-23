package logx

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

type HttpWriter struct {
	url string
}

func NewHttpWriter(url string) *HttpWriter {
	return &HttpWriter{url: url}
}

func (w *HttpWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	// 去掉末尾换行符
	body := bytes.NewBuffer(p[:len(p)-1])
	req, err := http.NewRequest("POST", w.url, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return len(p), nil
}
