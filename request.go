package curl

import (
	"bytes"
	"errors"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type Curl interface {
	// 設定檔案上傳
	SetFileData(fileData map[string]interface{}, filename string) error
	// 使用 Get Method 向對方請求
	Get() (resp *Response, err error)
	// 使用 Post Method 向對方請求
	Post() (resp *Response, err error)
	// 使用 Put Method 向對方請求
	Put() (resp *Response, err error)
	// 使用 Delete Method 向對方請求
	Delete() (resp *Response, err error)
	// 使用 Patch Method 向對方請求
	Patch() (resp *Response, err error)
	// SetWithURL TODO: v2 設定請求網址
	SetWithURL(url string) *Request
	// SetWithHeader TODO: v2 設定請求 headers
	SetWithHeader(headers map[string]string) *Request
	// SetWithCookie TODO: v2 設定請求 cookies
	SetWithCookie(cookies map[string]string) *Request
	// SetWithQuery TODO: v2 設定 URL 提供的參數,將用戶自定義的 url 參數添加到 http.Request 實例上
	SetWithQuery(queries map[string]interface{}) *Request
	// SetWithRawData TODO: v2 設定用戶提供的 raw data
	SetWithRawData(rawData map[string]interface{}) *Request
	// SetWithFormData TODO: v2 設定用戶提供的 form data
	SetWithFormData(formData map[string]interface{}) *Request
	// SetWithLock TODO: 發起 http ConnectPool 連線請求前需先上鎖
	SetWithLock() *Request
	// SetWithUnlock TODO: 結束 http ConnectPool 請求後需解鎖
	SetWithUnlock()
}

type Request struct {
	url       string                 // 請求網址
	headers   map[string]string      // 表頭內容
	cookies   map[string]string      // cookie 內容
	queries   map[string]interface{} // 網址參數內容
	method    string                 // 請求方式
	expire    time.Duration          // 連線*超時時間設定
	transport *http.Transport        // transport 規則
	skipTLS   bool                   // 憑證驗證狀態
	payload   io.Reader              // http request body 內容
	cli       *http.Client           // http.Client
	req       *http.Request          // http.Request 內容
	mutex     *sync.Mutex            // TODO: 互斥鎖用於ConnectPool連線模式
}

// NewRequest 創建一個 Request 實例
func NewRequest(options ...func(*Request)) Curl {
	// 初始化
	r := &Request{
		cli:       &http.Client{},
		req:       &http.Request{},
		payload:   nil,
		expire:    time.Second * 60,
		skipTLS:   true,
		transport: &http.Transport{},
		mutex:     &sync.Mutex{},
	}

	// running func list
	for _, o := range options {
		o(r)
	}

	return r
}

// SetClientRule 設定 cli 規則(ex: ssl 驗證是否跳過、timeout)
func (r *Request) setClientRule() {
	r.cli = &http.Client{
		Timeout:   r.expire,
		Transport: r.transport,
	}
}

/*
SetSkipTLSVerify 設定 ssl 憑證驗證，預設為 true
*/
func SetSkipTLSVerify(skipTLS bool) func(r *Request) {
	return func(r *Request) {
		r.skipTLS = skipTLS
	}
}

/*
SetTimeOut 設定 expire 連線過期時間，預設為 60 秒
*/
func SetTimeOut(expire time.Duration) func(r *Request) {
	return func(r *Request) {
		r.expire = expire
	}
}

func SetTransport(t *http.Transport) func(r *Request) {
	return func(r *Request) {
		r.transport = t
	}
}

/*
SetURL 設定請求網址
*/
func SetURL(url string) func(r *Request) {
	return func(r *Request) {
		r.url = url
	}
}

/*
SetHeader 設定請求 headers
*/
func SetHeader(headers map[string]string) func(r *Request) {
	return func(r *Request) {
		r.headers = headers
	}
}

// setHeader 將用戶自定義的表頭，添加到 http.Request 實例上
func (r *Request) setHeader() {
	for k := range r.headers {
		r.req.Header.Add(k, r.headers[k])
	}
}

/*
SetCookie 設定請求 cookies
*/
func SetCookie(cookies map[string]string) func(r *Request) {
	return func(r *Request) {
		r.cookies = cookies
	}
}

// setCookie 將用戶自定義的 cookies，添加到 http.Request 實例上
func (r *Request) setCookie() {
	for k := range r.cookies {
		r.req.AddCookie(&http.Cookie{
			Name:  k,
			Value: r.cookies[k],
		})
	}
}

/*
SetQuery 設定 URL 提供的參數,將用戶自定義的 url 參數添加到 http.Request 實例上
*/
func SetQuery(queries map[string]interface{}) func(r *Request) {
	return func(r *Request) {
		r.queries = queries
	}
}

func (r *Request) setQuery() {
	q := r.req.URL.Query()
	for k := range r.queries {
		paramV := reflect.ValueOf(r.queries[k])
		if paramV.Kind() == reflect.Slice {
			for i := 0; i < paramV.Len(); i++ {
				value := paramV.Index(i)
				q.Add(k, fmt.Sprint(value))
			}
			continue
		}
		q.Add(k, fmt.Sprint(paramV))
	}
	r.req.URL.RawQuery = q.Encode()
}

// SetRawData 設定用戶提供的 raw data
func SetRawData(rawData map[string]interface{}) func(r *Request) {

	return func(r *Request) {
		byteData, _ := jsoniter.Marshal(rawData)
		r.payload = strings.NewReader(string(byteData))
	}
}

/*
SetFormData 設定用戶提供的 form data
*/
func SetFormData(formData map[string]interface{}) func(r *Request) {

	return func(r *Request) {
		form := url.Values{}

		for k := range formData {
			paramV := reflect.ValueOf(formData[k])
			if paramV.Kind() == reflect.Slice {
				for i := 0; i < paramV.Len(); i++ {
					value := paramV.Index(i)
					form.Add(k, fmt.Sprint(value))
				}
				continue
			}
			form.Add(k, fmt.Sprint(paramV))
		}

		r.payload = strings.NewReader(form.Encode())
	}
}

func (r *Request) SetFileData(fileData map[string]interface{}, filename string) error {
	// the file data will be the second part of the body
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// use the writer to write the Part headers to the buffer
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// 組參數
	for pk, pv := range fileData {
		paramV := reflect.ValueOf(pv)
		if paramV.Kind() == reflect.Slice {
			for i := 0; i < paramV.Len(); i++ {
				value := paramV.Index(i)
				_ = writer.WriteField(pk, fmt.Sprint(value))
			}
			continue
		}
		_ = writer.WriteField(pk, fmt.Sprint(paramV))
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	// 設定 header
	r.headers["Content-Type"] = writer.FormDataContentType()

	r.payload = body

	return nil
}

// Get 使用 Get Method 向對方請求
func (r *Request) Get() (resp *Response, err error) {
	r.method = http.MethodGet
	return r.send()
}

// Post 使用 Post Method 向對方請求
func (r *Request) Post() (resp *Response, err error) {
	r.method = http.MethodPost
	return r.send()
}

// Put 使用 Put Method 向對方請求
func (r *Request) Put() (resp *Response, err error) {
	r.method = http.MethodPut
	return r.send()
}

// Delete 使用 Delete Method 向對方請求
func (r *Request) Delete() (resp *Response, err error) {
	r.method = http.MethodDelete
	return r.send()
}

// Patch 使用 Patch Method 向對方請求
func (r *Request) Patch() (resp *Response, err error) {
	r.method = http.MethodPatch
	return r.send()
}

func (r *Request) PostFile() (resp *Response, err error) {
	r.method = http.MethodPost

	return r.send()
}

// Send 發送 Http 請求
func (r *Request) send() (*Response, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover found：%s\n", r)
			debug.PrintStack()
		}
	}()

	// 檢查是否有 URL
	if r.url == "" {
		return nil, errors.New("Lock of request url")
	}

	// 檢查是否存在 Method
	if r.method == "" {
		return nil, errors.New("Lock of request method")
	}

	// 初始化 Response 對象
	resp := NewResponse()

	// 建立一個 http 請求
	var err error
	if r.req, err = http.NewRequest(r.method, r.url, r.payload); err != nil {
		return nil, err
	}

	// 設定 header 到 req
	r.setHeader()

	// 設定 cookie 到 req
	r.setCookie()

	// 設定 query 到 req
	r.setQuery()

	// 設定 client rule 到 req
	r.setClientRule()

	// 執行 API 請求
	if resp.Raw, err = r.cli.Do(r.req); err != nil {
		return nil, err
	}
	defer resp.Raw.Body.Close()

	// check response status
	if !resp.isOk() {
		return resp, errors.New("Lock of request method")
	}

	// 處理 http response header
	resp.parseHeaders()

	// 處理 http response body
	if err = resp.parseBody(); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Request) SetWithLock() *Request {
	r.mutex.Lock()
	return r
}

func (r *Request) SetWithUnlock() {
	r.mutex.Unlock()
}

func (r *Request) SetWithURL(url string) *Request {
	r.url = url
	return r
}

func (r *Request) SetWithHeader(headers map[string]string) *Request {
	r.headers = headers
	return r
}

func (r *Request) SetWithCookie(cookies map[string]string) *Request {
	r.cookies = cookies
	return r
}

func (r *Request) SetWithQuery(queries map[string]interface{}) *Request {
	r.queries = queries
	return r
}

func (r *Request) SetWithFormData(formData map[string]interface{}) *Request {
	form := url.Values{}

	for k := range formData {
		paramV := reflect.ValueOf(formData[k])
		if paramV.Kind() == reflect.Slice {
			for i := 0; i < paramV.Len(); i++ {
				value := paramV.Index(i)
				form.Add(k, fmt.Sprint(value))
			}
			continue
		}
		form.Add(k, fmt.Sprint(paramV))
	}

	r.payload = strings.NewReader(form.Encode())
	return r
}

func (r *Request) SetWithRawData(rawData map[string]interface{}) *Request {
	byteData, _ := jsoniter.Marshal(rawData)
	r.payload = strings.NewReader(string(byteData))
	return r
}
