package bybit

import (
	"bytes"
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
	"strconv"
	"strings"
	"time"
)

const (
	// MainNetBaseURL :
	MainNetBaseURL = "https://api.bybit.com"
	// MainNetBaseURL2 :
	MainNetBaseURL2 = "https://api.bytick.com"
)

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

// Client :
type Client struct {
	httpClient *http.Client

	baseURL string
	key     string
	secret  string

	debug  bool
	logger Logger

	checkResponseBody checkResponseBodyFunc
}

// NewClient :
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},

		baseURL:           MainNetBaseURL,
		checkResponseBody: checkResponseBody,
	}
}

// WithHTTPClient :
func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient

	return c
}

// WithHTTPClient :
func (c *Client) Debug(logger Logger) *Client {
	c.debug = true
	c.logger = logger

	return c
}

// WithAuth :
func (c *Client) WithAuth(key string, secret string) *Client {
	c.key = key
	c.secret = secret

	return c
}

func (c Client) withCheckResponseBody(f checkResponseBodyFunc) *Client {
	c.checkResponseBody = f

	return &c
}

// WithBaseURL :
func (c *Client) WithBaseURL(url string) *Client {
	c.baseURL = url

	return c
}

type RateLimitHeaders struct {
	RateLimitStatus  int `json:"rate_limit_status"`
	RateLimitResetMs int `json:"rate_limit_reset_ms"`
	RateLimit        int `json:"rate_limit"`
}

// Request :
func (c *Client) V5Request(req *http.Request, dst interface{}) error {

	if c.debug {
		c.logger.Debugf("Request url: %s", req.URL.String())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.debug {
		c.logger.Debugf("Response: %v", resp)
	}

	headers := RateLimitHeaders{}

	switch {
	case 200 <= resp.StatusCode && resp.StatusCode <= 299:
		body, err := io.ReadAll(resp.Body)

		if c.debug {
			c.logger.Debugf("Body: %s", string(body))
		}

		if err != nil {
			if c.debug {
				c.logger.Errorf("Error: %s", err.Error())
			}
			return err
		}

		if c.checkResponseBody == nil {
			return errors.New("checkResponseBody func should be set")
		}
		if err := c.checkResponseBody(body); err != nil {
			return err
		}

		if err := json.Unmarshal(body, &dst); err != nil {
			return err
		}

		// headers.RateLimitStatus =
		if v, err := strconv.Atoi(resp.Request.Header.Get("X-Bapi-Limit-Status")); err == nil {
			headers.RateLimitStatus = v
		}

		if v, err := strconv.Atoi(resp.Request.Header.Get("X-Bapi-Limit-Reset-Timestamp")); err == nil {
			headers.RateLimitResetMs = v
		}

		if v, err := strconv.Atoi(resp.Request.Header.Get("X-Bapi-Limit")); err == nil {
			headers.RateLimit = v
		}

		if h, err := json.Marshal(headers); err == nil {
			if err := json.Unmarshal(h, &dst); err != nil {
				return err
			}
		}

		return nil
	case resp.StatusCode == http.StatusForbidden:
		return ErrAccessDenied
	case resp.StatusCode == http.StatusNotFound:
		return ErrPathNotFound
	default:

		body, err := io.ReadAll(resp.Body)
		if c.debug {
			c.logger.Debugf("Body: %v", string(body))
		}

		if err != nil && c.debug {
			c.logger.Errorf("Error: %s", err.Error())
		}

		return errors.New("unexpected error")
	}
}

// Request :
func (c *Client) Request(req *http.Request, dst interface{}) error {

	if c.debug {
		c.logger.Debugf("Request url: %s", req.URL.String())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.debug {
		c.logger.Debugf("Response: %v", resp)
	}

	switch {
	case 200 <= resp.StatusCode && resp.StatusCode <= 299:
		body, err := io.ReadAll(resp.Body)
		if c.debug {
			c.logger.Debugf("Body: %s", string(body))
		}

		if err != nil {
			if c.debug {
				c.logger.Errorf("Error: %s", err.Error())
			}
			return err
		}

		if c.checkResponseBody == nil {
			return errors.New("checkResponseBody func should be set")
		}
		if err := c.checkResponseBody(body); err != nil {
			return err
		}

		if err := json.Unmarshal(body, &dst); err != nil {
			return err
		}
		return nil
	case resp.StatusCode == http.StatusForbidden:
		if c.debug {
			c.logger.Errorf("Error: %d", http.StatusForbidden)
		}
		return ErrAccessDenied
	case resp.StatusCode == http.StatusNotFound:
		if c.debug {
			c.logger.Errorf("Error: %d", http.StatusNotFound)
		}
		return ErrPathNotFound
	default:

		body, err := io.ReadAll(resp.Body)
		if c.debug {
			c.logger.Debugf("Body: %v", string(body))
		}

		if err != nil && c.debug {
			c.logger.Errorf("Error: %s", err.Error())
		}

		return errors.New("unexpected error")
	}
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

func (c *Client) populateSignatureForBody(src []byte) []byte {
	intNow := int(time.Now().UTC().UnixNano() / int64(time.Millisecond))
	now := strconv.Itoa(intNow)

	body := map[string]interface{}{}
	if err := json.Unmarshal(src, &body); err != nil {
		panic(err)
	}

	body["api_key"] = c.key
	body["timestamp"] = now
	body["sign"] = getSignatureForBody(body, c.secret)

	result, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	return result
}

func getV5Signature(
	timestamp int,
	key string,
	queryString string,
	secret string,
) string {
	val := strconv.Itoa(timestamp) + key
	val = val + queryString
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(val))
	return hex.EncodeToString(h.Sum(nil))
}

func getV5SignatureForBody(
	timestamp int,
	key string,
	body []byte,
	secret string,
) string {
	val := strconv.Itoa(timestamp) + key
	val = val + string(body)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(val))
	return hex.EncodeToString(h.Sum(nil))
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

func getSignatureForBody(src map[string]interface{}, key string) string {
	keys := make([]string, len(src))
	i := 0
	_val := ""
	for k := range src {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		_val += k + "=" + fmt.Sprintf("%v", src[k]) + "&"
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

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	if err := c.Request(req, &dst); err != nil {
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

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	if err := c.Request(req, &dst); err != nil {
		return err
	}

	return nil
}

func (c *Client) getV5Privately(path string, query url.Values, dst interface{}) error {
	if !c.hasAuth() {
		return fmt.Errorf("this is private endpoint, please set api key and secret")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	u.Path = path
	u.RawQuery = query.Encode()

	timestamp := int(time.Now().UTC().UnixNano() / int64(time.Millisecond))
	sign := getV5Signature(timestamp, c.key, query.Encode(), c.secret)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-BAPI-API-KEY", c.key)
	req.Header.Set("X-BAPI-TIMESTAMP", strconv.Itoa(timestamp))
	req.Header.Set("X-BAPI-SIGN", sign)

	if err := c.V5Request(req, &dst); err != nil {
		return err
	}

	return nil
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

	body = c.populateSignatureForBody(body)

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if err := c.Request(req, &dst); err != nil {
		return err
	}

	return nil
}

func (c *Client) postV5JSON(path string, body []byte, dst interface{}) error {
	if !c.hasAuth() {
		return fmt.Errorf("this is private endpoint, please set api key and secret")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	u.Path = path

	timestamp := int(time.Now().UTC().UnixNano() / int64(time.Millisecond))
	sign := getV5SignatureForBody(timestamp, c.key, body, c.secret)

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BAPI-API-KEY", c.key)
	req.Header.Set("X-BAPI-TIMESTAMP", strconv.Itoa(timestamp))
	req.Header.Set("X-BAPI-SIGN", sign)

	if err := c.Request(req, &dst); err != nil {
		return err
	}

	return nil
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

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	if err := c.Request(req, &dst); err != nil {
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

	if err := c.Request(req, &dst); err != nil {
		return err
	}

	return nil
}
