package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"gitee.com/lihp1603/utils/log"
)

type HttpApiClient struct {
	// HTTP client used to communicate with the API. By default
	// http.DefaultClient will be used.
	client *http.Client
}

// NewClient uses an http.Transport with reasonable defaults.
func NewHttpApiClient() *HttpApiClient {
	return NewHttpApiClientWithConfig(nil)
}

// NewClient uses an http.Transport with reasonable defaults.
func NewHttpsApiClient(cert tls.Certificate) *HttpApiClient {
	return NewHttpApiClientWithConfig(&tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
	})
}

// Therefore, the config.Certificates must contain a TLS
// certificate that is valid for client authentication.
//
// NewClientWithConfig uses an http.Transport with reasonable
// defaults.
func NewHttpApiClientWithConfig(config *tls.Config) *HttpApiClient {
	return &HttpApiClient{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second, //限制建立TCP连接的时间
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second, //限制tls握手时间
				ResponseHeaderTimeout: 10 * time.Second, //限制读取response header的时间
				ExpectContinueTimeout: 1 * time.Second,  //限制client在发送包含expect:100-continue的header到收到继续发送body的response之间的时间等待
				TLSClientConfig:       config,
			},
		},
	}
}

// Send creates a new HTTP request with the given method, context
// request body and request options, if any. It randomly iterates
// over the given endpoints until it receives a HTTP response.
//
// If sending a request to one endpoint fails due to e.g. a network
// or DNS error, Send tries the next endpoint. It aborts once the
// context is canceled or its deadline exceeded.
func (c *HttpApiClient) Send(ctx context.Context, method string, endpoints []string, path string, body io.ReadSeeker) (*http.Response, error) {
	if len(endpoints) == 0 {
		return nil, errors.New("no server endpoint")
	}
	var (
		request  *http.Request
		response *http.Response
		err      error
		R        = rand.Intn(len(endpoints)) // randomize endpoints => avoid hitting the same endpoint all the time.
	)
	for i := range endpoints {
		nextEndpoint := endpoints[(i+R)%len(endpoints)]
		request, err = http.NewRequestWithContext(ctx, method, endpoint(nextEndpoint, path), retryBody(body))
		if err != nil {
			return nil, err
		}
		//for _, opt := range options {
		//	opt(request)
		//}

		response, err = c.Do(request)
		if err == nil {
			return response, nil
		}
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
	}
	return response, err
}

// Get issues a GET to the specified URL.
// It is a wrapper around retry.Do.
func (c *HttpApiClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, retryBody(nil))
	if err != nil {
		log.Error("http get:%s err:%s", url, err.Error())
		return nil, err
	}
	return c.Do(req)
}

// Post issues a POST to the specified URL.
// It is a wrapper around retry.Do.
func (c *HttpApiClient) Post(ctx context.Context, url, contentType string, body io.ReadSeeker) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, retryBody(body))
	if err != nil {
		log.Error("http post:%s err:%s", url, err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

// Do sends an HTTP request and returns an HTTP response using
// the underlying http.Client. If the request fails b/c of a
// temporary error Do retries the request a few times. If the
// request keeps failing, Do will give up and return a descriptive
// error.
func (c *HttpApiClient) Do(req *http.Request) (*http.Response, error) {
	type RetryReader interface {
		io.Reader
		io.Seeker
		io.Closer
	}

	// If the request body is not a RetryReader it cannot
	// be retried. The caller has to ensure that the actual
	// body content is an io.ReadCloser + io.Seeker.
	// The retry.NewRequest method does that.
	//
	// A request can only be retried if we can seek to the
	// start of the request body. Otherwise, we may send a
	// partial response body when we retry the request.
	var body RetryReader
	if req.Body != nil {
		var ok bool
		body, ok = req.Body.(RetryReader)
		if !ok {
			// We cannot convert the req.Body to an io.Seeker.
			// If we would proceed we may introduce hard to find
			// bugs. Also, there is no point in returning an
			// error since the caller has specified a wrong type.
			panic("request cannot be retried")
		}

		// If there is a request body, additionally set the
		// GetBody callback - if not set already. The underlying
		// HTTP stack will use the GetBody callback to obtain a new
		// copy of the request body - e.g. in case of a redirect.
		if req.GetBody == nil {
			req.GetBody = func() (io.ReadCloser, error) {
				if _, err := body.Seek(0, io.SeekStart); err != nil {
					return nil, err
				}
				return body, nil
			}
		}
	}

	const (
		MinRetryDelay     = 200 * time.Millisecond
		MaxRandRetryDelay = 800
	)
	var (
		retry  = 2 // For now, we retry 2 times before we give up
		client = c.client
	)
	resp, err := client.Do(req)
	for retry > 0 && (isTemporary(err) || (resp != nil && resp.StatusCode == http.StatusServiceUnavailable)) {
		randomRetryDelay := time.Duration(rand.Intn(MaxRandRetryDelay)) * time.Millisecond
		time.Sleep(MinRetryDelay + randomRetryDelay)
		retry--

		// If there is a body we have to reset it. Otherwise, we may send
		// only partial data to the server when we retry the request.
		if body != nil {
			if _, err = body.Seek(0, io.SeekStart); err != nil {
				return nil, err
			}
			req.Body = body
		}

		resp, err = client.Do(req) // Now, retry.
	}
	if isTemporary(err) {
		// If the request still fails with a temporary error
		// we wrap the error to provide more information to the
		// caller.
		return nil, &url.Error{
			Op:  req.Method,
			URL: req.URL.String(),
			Err: fmt.Errorf("Temporary network error: %v", err),
		}
	}
	return resp, err
}

// endpoint returns an endpoint URL starting with the
// given endpoint followed by the path elements.
//
// For example:
//   • endpoint("https://127.0.0.1:7373", "version")                => "https://127.0.0.1:7373/version"
//   • endpoint("https://127.0.0.1:7373/", "/key/create", "my-key") => "https://127.0.0.1:7373/key/create/my-key"
//
// Any leading or trailing whitespaces are removed from
// the endpoint before it is concatenated with the path
// elements.
//
// The path elements will not be URL-escaped.
func endpoint(endpoint string, elems ...string) string {
	endpoint = strings.TrimSpace(endpoint)
	endpoint = strings.TrimSuffix(endpoint, "/")

	if len(elems) > 0 && !strings.HasPrefix(elems[0], "/") {
		endpoint += "/"
	}
	return endpoint + path.Join(elems...)
}

// isTemporary returns true if the given error is
// temporary - e.g. a temporary *url.Error or an
// net.Error that indicates that a request got
// timed-out.
//
// A nil error is not temporary.
func isTemporary(err error) bool {
	if err == nil { // fast path
		return false
	}
	if netErr, ok := err.(net.Error); ok { // *url.Error implements net.Error
		if netErr.Timeout() || netErr.Temporary() {
			return true
		}

		// If a connection drops (e.g. server dies) while sending the request
		// http.Do returns either io.EOF or io.ErrUnexpected. We treat that as
		// temp. since the server may get restared such that the retry may succeed.
		if errors.Is(netErr, io.EOF) || errors.Is(netErr, io.ErrUnexpectedEOF) {
			return true
		}
	}
	return false
}

// retryBody takes an io.ReadSeeker and converts it
// into an io.ReadCloser that can be used as request
// body for retryable requests.
//
// The body must implement io.Seeker to ensure that
// the entire body is sent again when retrying a request.
//
// If body is nil, retryBody returns nil.
func retryBody(body io.ReadSeeker) io.ReadCloser {
	if body == nil {
		return nil
	}

	var closer io.Closer
	if c, ok := body.(io.Closer); ok {
		closer = c
	} else {
		closer = ioutil.NopCloser(body)
	}

	type ReadSeekCloser struct {
		io.ReadSeeker
		io.Closer
	}
	return ReadSeekCloser{
		ReadSeeker: body,
		Closer:     closer,
	}
}
