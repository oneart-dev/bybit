package bybit

import "os"

const (
	// TestNetBaseURL :
	TestNetBaseURL = "https://api-testnet.bybit.com"
)

// TestClient :
type TestClient struct {
	*Client
}

// NewTestClient :
func NewTestClient() *TestClient {
	return &TestClient{
		Client: &Client{
			baseURL: TestNetBaseURL,
		},
	}
}

// WithAuthFromEnv :
func (c *TestClient) WithAuthFromEnv() *TestClient {
	key, ok := os.LookupEnv("BYBIT_TEST_KEY")
	if !ok {
		panic("need BYBIT_TEST_KEY as environment variable")
	}
	secret, ok := os.LookupEnv("BYBIT_TEST_SECRET")
	if !ok {
		panic("need BYBIT_TEST_SECRET as environment variable")
	}
	c.key = key
	c.secret = secret

	return c
}

// WithBaseURL :
func (c *TestClient) WithBaseURL(url string) *TestClient {
	c.baseURL = url

	return c
}
