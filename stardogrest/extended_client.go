package stardogrest

import (
	"crypto/tls"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type ExtendedBaseClient struct {
	BaseClient
}

func (client *ExtendedBaseClient) SetConnection(url, username, password string) {
	client.RetryDuration = 3 * time.Second
	client.RetryAttempts = 1
	client.BaseClient.Sender = createCustomSender()
	client.BaseClient.BaseURI = getFullURL(url)
	client.BaseClient.Authorizer = autorest.NewBasicAuthorizer(username, password)
}

func NewExtendedBaseClient() *ExtendedBaseClient {
	return &ExtendedBaseClient{
		BaseClient: New(),
	}
}

func getFullURL(incompleteURL string) string {
	uri, _ := url.ParseRequestURI(incompleteURL)
	uri.Path = "/admin"
	return uri.Redacted()
}

func createCustomSender() autorest.Sender {
	defaultTransport := http.DefaultTransport.(*http.Transport)
	transport := &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSClientConfig: &tls.Config{
			MinVersion:    tls.VersionTLS12,
			Renegotiation: tls.RenegotiateNever,
		},
	}
	var roundTripper http.RoundTripper = transport
	if tracing.IsEnabled() {
		roundTripper = tracing.NewTransport(transport)
	}
	j, _ := cookiejar.New(nil)
	return &http.Client{Jar: j, Transport: roundTripper, Timeout: time.Second}
}
