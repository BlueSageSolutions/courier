package transform

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BlueSageSolutions/courier/pkg/util"
	"golang.org/x/oauth2"
)

const (
	skipTLSVerify      string = "SKIP_TLS_VERIFY"
	tlsCertificatePath string = "TLS_CERTIFICATE_PATH"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: time.Minute * 30,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			Proxy:               http.ProxyFromEnvironment,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig:     getTLSConfig(),
		},
	}
}

func getTLSConfig() *tls.Config {
	var tlsCfg tls.Config
	skipTLS := os.Getenv(skipTLSVerify)
	tlsCfg.InsecureSkipVerify = strings.EqualFold(skipTLS, "true")

	cert := os.Getenv(tlsCertificatePath)
	if cert == "" {
		return &tlsCfg
	}
	caCert, err := os.ReadFile(cert)
	if err != nil {
		util.GetLogger().Info(fmt.Sprintf("unable to find client certificate: %s", tlsCertificatePath))

		return &tlsCfg
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsCfg.RootCAs = caCertPool
	return &tlsCfg
}

// AddTLSConfigToTransport adds TLS Config to http.RoundTripper, according to few environment variables.
func AddTLSConfigAndProxyToTransport(tr http.RoundTripper) {
	if tr == nil {
		return
	}

	if tr, ok := tr.(*oauth2.Transport); ok {
		if tr.Base == nil {
			tr.Base = &http.Transport{
				Proxy:           http.ProxyFromEnvironment,
				TLSClientConfig: getTLSConfig(),
			}
			return
		}

		if tk, ok := tr.Base.(*http.Transport); ok {
			tk.Proxy = http.ProxyFromEnvironment
			tk.TLSClientConfig = getTLSConfig()
			tr.Base = tk
		}
	}
}

// Do do the request using the custom http client.
func Do(req *http.Request) (*http.Response, error) {
	return client.Do(req)
}

// GetResponse get the response object from the request.
func GetResponse(req *http.Request) ([]byte, error) {
	resp, err := Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	return d, err
}

// GetResponseWithStatus makes the request, read the body as bytes and status code from the response.
func GetResponseWithStatus(req *http.Request) (int, []byte, error) {
	resp, err := Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	return resp.StatusCode, d, err
}

// ParseJSON get the JSON response body.
func ParseJSON(req *http.Request, v interface{}) error {
	_, err := ParseJSONWithStatus(req, v)
	return err
}

// ParseJSONWithStatus make the request, unmarshals the response body into the second argument and return the status code.
func ParseJSONWithStatus(req *http.Request, v interface{}) (int, error) {
	status, d, err := GetResponseWithStatus(req)
	if err != nil {
		return status, err
	}
	err = json.Unmarshal(d, v)
	if err != nil {
		return status, err
	}
	return status, nil
}

// ParseJSONWithResponseHeader get the JSON response body and header status.
func ParseJSONWithResponseHeader(req *http.Request, v interface{}) (http.Header, error) {
	headers, d, err := getResponseWithHeaders(req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(d, v)
	if err != nil {
		return headers, err
	}
	return headers, nil
}

// getResponseWithHeaders get the response object from the request.
func getResponseWithHeaders(req *http.Request) (http.Header, []byte, error) {
	resp, err := Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	return resp.Header, d, err
}
