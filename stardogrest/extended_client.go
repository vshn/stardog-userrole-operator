package stardogrest

import (
	"crypto/tls"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/tracing"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

type ExtendedBaseClient struct {
	stardog client.StardogClient
	BaseClient
}

func (c *ExtendedBaseClient) SetConnection(url, username, password string) {
	test := client.New(httptransport.New("", "", nil), strfmt.Default)
	cred := httptransport.BasicAuth(os.Getenv("API_ACCESS_TOKEN"), "a")

	test.Roles.ListRoles(&roles.ListRolesParams{}, cred)
	test.Users.ChangePassword()
	c.RetryDuration = 3 * time.Second
	c.RetryAttempts = 1
	c.BaseClient.Sender = createCustomSender()
	c.BaseClient.BaseURI = getFullURL(url)
	c.BaseClient.Authorizer = autorest.NewBasicAuthorizer(username, password)
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
