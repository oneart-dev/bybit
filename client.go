package bybit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// MainNetBaseURL :
	MainNetBaseURL = "https://api.bybit.com"
)

// Client :
type Client struct {
	baseURL    string
	key        string
	secret     string
	httpClient *http.Client
	debug      bool
}

// NewClient :
func NewClient() *Client {
	return &Client{
		baseURL:    MainNetBaseURL,
		httpClient: &http.Client{},
	}
}

// SetDebug :
func (c *Client) Debug() *Client {
	c.debug = true
	return c
}

// WithHTTPClient :
func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

// WithAuth :
func (c *Client) WithAuth(key string, secret string) *Client {
	c.key = key
	c.secret = secret

	return c
}

// hasAuth : check has auth key and secret
func (c *Client) hasAuth() bool {
	return c.key != "" && c.secret != ""
}

func (c *Client) populateSignature(src url.Values) url.Values {
	intNow := int(time.Now().UTC().UnixNano() / int64(time.Millisecond))
	now := strconv.Itoa(intNow)

	if src == nil {
		src = url.Values{}
	}

	src.Add("api_key", c.key)
	src.Add("timestamp", now)
	src.Add("sign", getSignature(src, c.secret))

	return src
}

func getSignature(src url.Values, key string) string {
	keys := make([]string, len(src))
	i := 0
	_val := ""
	for k := range src {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		_val += k + "=" + src.Get(k) + "&"
	}
	_val = _val[0 : len(_val)-1]
	h := hmac.New(sha256.New, []byte(key))
	_, err := io.WriteString(h, _val)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *Client) getPublicly(path string, query url.Values, dst interface{}) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	u.Path = path
	u.RawQuery = query.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.debugResponse(resp)

	if c.checkResponseForErrors(resp) != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&dst); err != nil {
		return err
	}

	return nil
}

func (c *Client) getPrivately(path string, query url.Values, dst interface{}) error {
	if !c.hasAuth() {
		return fmt.Errorf("this is private endpoint, please set api key and secret")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	u.Path = path
	query = c.populateSignature(query)
	u.RawQuery = query.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.checkResponseForErrors(resp) != nil {
		return err
	}

	c.debugResponse(resp)

	if c.checkResponseForErrors(resp) != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&dst); err != nil {
		return err
	}

	return nil
}

var ErrTooManyRequests = errors.New("too many requests")
var ErrServerFail = errors.New("server fail")

func (c *Client) checkResponseForErrors(resp *http.Response) error {
	if resp.StatusCode == 200 {
		return nil
	}

	if resp.StatusCode == 429 || resp.StatusCode == 403 {
		return ErrTooManyRequests
	}

	if resp.StatusCode >= 500 {
		return ErrServerFail
	}

	return errors.New("request failed: " + resp.Status)
}

func (c *Client) postJSON(path string, body []byte, dst interface{}) error {
	if !c.hasAuth() {
		return fmt.Errorf("this is private endpoint, please set api key and secret")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	u.Path = path

	query := url.Values{}
	query = c.populateSignature(query)
	u.RawQuery = query.Encode()

	resp, err := c.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.debugResponse(resp)

	if c.checkResponseForErrors(resp) != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&dst); err != nil {
		return err
	}

	return nil
}

func (c *Client) debugResponse(resp *http.Response) {
	if c.debug {
		fmt.Println("RESPONSE DEBUG INFORMATION")
		fmt.Println("Query: ", resp.Request.URL.String())
		fmt.Println("Status: ", resp.Status)
		fmt.Println("Headers: ", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Body: ", string(body))
		resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
}

func (c *Client) postForm(path string, body url.Values, dst interface{}) error {
	if !c.hasAuth() {
		return fmt.Errorf("this is private endpoint, please set api key and secret")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil
	}
	u.Path = path

	body = c.populateSignature(body)

	resp, err := c.httpClient.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.debugResponse(resp)

	if c.checkResponseForErrors(resp) != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&dst); err != nil {
		return err
	}

	return nil
}

func (c *Client) deletePrivately(path string, query url.Values, dst interface{}) error {
	if !c.hasAuth() {
		return fmt.Errorf("this is private endpoint, please set api key and secret")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	u.Path = path
	query = c.populateSignature(query)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.debugResponse(resp)

	if c.checkResponseForErrors(resp) != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&dst); err != nil {
		return err
	}

	return nil
}
