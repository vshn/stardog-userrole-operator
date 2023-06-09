package stardogrest

import (
	"crypto/tls"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/tracing"
	client2 "github.com/go-openapi/runtime/client"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type ExtendedBaseClient struct {
	client client.StardogClient
}

func (c *ExtendedBaseClient) SetConnection(url, username, password string) {
	c.client.
	config := client.DefaultTransportConfig().WithHost(url)
	cl := client.NewHTTPClientWithConfig(nil, config)
	auth := client2.BasicAuth(username, password)
	cl.Users.ChangePassword(nil, auth)
}

func NewExtendedBaseClient() *ExtendedBaseClient {
	return &ExtendedBaseClient{
		client: *client.NewHTTPClient(nil),
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
