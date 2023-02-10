package service

import (
	curl "go-http-curl"
	"net/http"
	"sync"
	"time"
)

type Interface interface {
	MethodA()
	MethodB()
	MethodC()
	MethodD()
	MethodF()
}

type service struct {
	curl curl.Curl
}

var singleton *service
var once sync.Once

func NewService() Interface {
	once.Do(func() {
		singleton = &service{
			curl: curl.NewRequest(curl.SetTransport(
				&http.Transport{
					MaxIdleConns:          100,
					MaxIdleConnsPerHost:   100,
					MaxConnsPerHost:       100,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
				})),
		}
	})
	return singleton
}

// MethodA 對127.0.0.1:8087發出請求
func (s *service) MethodA() {
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Bearer aaaa",
	}
	params := map[string]interface{}{
		"lang": "zh-tw",
	}

	s.curl.SetWithLock().
		SetWithURL("http://127.0.0.1:8087").
		SetWithHeader(headers).
		SetWithQuery(params).
		Get()

	s.curl.SetWithUnlock()
}

// MethodB 對127.0.0.1:8088發出請求
func (s *service) MethodB() {
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Bearer bbbb",
	}
	params := map[string]interface{}{
		"lang": "en",
	}

	s.curl.SetWithLock().
		SetWithURL("http://127.0.0.1:8088").
		SetWithHeader(headers).
		SetWithQuery(params).
		Get()

	s.curl.SetWithUnlock()
}

// MethodC 對127.0.0.1:8089發出請求
func (s *service) MethodC() {
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Bearer cccc",
	}
	params := map[string]interface{}{
		"lang": "jp",
	}

	s.curl.SetWithLock().
		SetWithURL("http://127.0.0.1:8089").
		SetWithHeader(headers).
		SetWithQuery(params).
		Get()

	s.curl.SetWithUnlock()
}

func (s *service) MethodD() {
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Bearer dddd",
	}
	params := map[string]interface{}{
		"lang": "zh-tw",
	}

	req := curl.NewRequest()

	req.SetWithURL("http://127.0.0.1:8087").
		SetWithHeader(headers).
		SetWithQuery(params).
		Get()
}

func (s *service) MethodF() {
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Bearer ffff",
	}
	params := map[string]interface{}{
		"lang": "en",
	}
	req := curl.NewRequest()

	req.SetWithURL("http://127.0.0.1:8088").
		SetWithHeader(headers).
		SetWithQuery(params).
		Get()

}
