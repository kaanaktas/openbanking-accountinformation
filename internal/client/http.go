package client

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SecureClient struct {
	*http.Client
	endpoint string
	http.Header
}

type HttpResponse struct {
	StatusCode int
	Body       string
	Header     http.Header
}

var cacheMem = cache.LoadInMemory()

func NewSecureHttpClient(endpoint string, header http.Header) (*SecureClient, error) {
	transport, err := newHttpTransport()
	if err != nil {
		return nil, err
	}

	client := &SecureClient{
		&http.Client{
			Transport: transport,
			//disable direction follow
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		endpoint,
		header,
	}

	return client, nil
}

func newHttpTransport() (*http.Transport, error) {
	if value, found := cacheMem.Get("http_transport"); found {
		return value.(*http.Transport), nil
	}

	serverCACertPEM := os.Getenv("CLIENT_CA_CERT_PEM")
	clientCertPEM := os.Getenv("CLIENT_CERT_PEM")
	clientKeyPEM := os.Getenv("CLIENT_KEY_PEM")

	pool := x509.NewCertPool()

	for _, cert := range strings.Split(serverCACertPEM, ",") {
		caCert, err := ioutil.ReadFile(cert)
		if err != nil {
			return nil, errors.WithMessage(err, "error in newHttpTransport() while reading cert file")
		}
		pool.AppendCertsFromPEM(caCert)
	}

	clientCert, err := tls.LoadX509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return nil, errors.WithMessage(err, "error in newHttpTransport() while loading x509 certificate")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{clientCert},
		}}

	_ = cacheMem.Set("http_transport", transport, cache.NoExpiration)

	return transport, nil
}

func (s *SecureClient) Post(payload io.Reader) (*HttpResponse, error) {
	req, err := s.createRequest("POST", payload)
	if err != nil {
		return nil, errors.WithMessage(err, "error in Post() while creating NewRequest")
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "error in Post()")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "error in Post() in while reading response body")
	}

	response := &HttpResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Header:     resp.Header,
	}

	return response, nil
}

func (s *SecureClient) Get(parameters url.Values) (*HttpResponse, error) {
	req, err := s.createRequest("GET", nil)
	if err != nil {
		return nil, errors.WithMessage(err, "error in Get() while creating NewRequest")
	}

	req.URL.RawQuery = encode(parameters)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "error in Get()")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "error in Post() while creating NewRequest")
	}

	response := &HttpResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Header:     resp.Header,
	}

	return response, nil
}

func (s *SecureClient) createRequest(method string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, s.endpoint, payload)
	if err != nil {
		return nil, errors.WithMessage(err, "error in createRequest()")
	}

	req.Header = s.Header
	return req, nil
}

//custom uri encode. only replaces "space" with +(plus)
func encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}

	for _, k := range keys {
		vs := v[k]
		keyEscaped := k
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(strings.ReplaceAll(v, " ", "+"))
		}
	}
	return buf.String()
}
