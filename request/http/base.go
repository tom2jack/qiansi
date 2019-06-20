package http

import (
	"crypto/tls"
	"net/http"
	"time"
)

var (
	ZMHTTP *http.Client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 15,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
)

const BASE_URL = "https://localhost:8001"
