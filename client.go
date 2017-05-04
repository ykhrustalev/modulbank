package modulbank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
)

const ApiUrl = "https://api.modulbank.ru"

type Options struct {
	HttpClient *http.Client
	Sandbox    bool
	Token      string
	Logger     *logrus.Entry
}

type Client struct {
	httpClient *http.Client
	sandbox    bool
	token      string
	logger     *logrus.Entry
}

func NewClient(options *Options) *Client {
	httpClient := options.HttpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	logger := options.Logger
	if logger == nil {
		logger = logrus.WithField("context", "modulbank")
	}

	return &Client{
		httpClient: httpClient,
		sandbox:    options.Sandbox,
		token:      options.Token,
		logger:     logger,
	}
}

func (client *Client) buildUrl(path string) string {
	return ApiUrl + path
}

func (client *Client) buildRequest(logger *logrus.Entry, method, path string, data interface{}) (*http.Request, error) {
	var reqData io.Reader

	if data == nil {
		reqData = nil
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			logger.WithError(err).Error("failed to marshal json")
			return nil, err
		}
		reqData = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, client.buildUrl(path), reqData)
	if err != nil {
		logger.WithError(err).Error("failed to build request")
		return nil, err
	}

	client.preProcessRequest(req)

	return req, nil
}

func (client *Client) preProcessRequest(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	if client.sandbox {
		req.Header.Set("Authorization", "Bearer sandboxtoken")
		req.Header.Set("sandbox", "on")
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.token))
	}
}

func (client *Client) doRequest(logger *logrus.Entry, req *http.Request, expectCode int) (*http.Response, error) {
	rawBytes, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logger.WithError(err).Error("failed to dump request")
		return nil, err
	}
	logger.WithField("contents", string(rawBytes)).Debug("raw request")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		logger.WithError(err).Error("failed to do request")
		return nil, err
	}

	rawBytes, err = httputil.DumpResponse(resp, true)
	if err != nil {
		logger.WithError(err).Error("failed to dump response")
		return nil, err
	}
	logger.WithField("contents", string(rawBytes)).Debug("raw response")

	if resp.StatusCode != expectCode {
		logger.WithField("statusCode", resp.StatusCode).Error("unexpected response")
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return resp, nil
}

func (client *Client) decodeBody(logger *logrus.Entry, body io.Reader, target interface{}) error {
	err := json.NewDecoder(body).Decode(&target)
	if err != nil {
		logger.WithError(err).Error("failed to decode body")
		return err
	}
	return nil
}

func (client *Client) closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

func (client *Client) handleRequest(
	logger *logrus.Entry,
	method, path string,
	reqData interface{},
	result interface{},
) (*http.Response, error) {

	req, err := client.buildRequest(logger, method, path, reqData)
	if err != nil {
		return nil, err
	}

	resp, err := client.doRequest(logger, req, 200)
	defer client.closeBody(resp)
	if err != nil {
		return nil, err
	}

	err = client.decodeBody(logger, resp.Body, &result)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
