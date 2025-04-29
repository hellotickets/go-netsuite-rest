package netsuite

import (
	"net/http"
	"net/url"

	"github.com/hellotickets/go-netsuite-rest/utils"
)

func (c *Client) NewCustomRequest() CustomRequest {
	r := CustomRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type CustomRequest struct {
	client      *Client
	queryParams *CustomRequestQueryParams
	pathParams  *CustomRequestPathParams
	method      string
	headers     http.Header
	requestBody CustomRequestBody
}

func (r CustomRequest) NewQueryParams() *CustomRequestQueryParams {
	return &CustomRequestQueryParams{}
}

type CustomRequestQueryParams struct {
	Script int `schema:"script,omitempty"`
	Deploy int `schema:"deploy,omitempty"`
}

func (p CustomRequestQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	encoder.RegisterEncoder(Date{}, utils.EncodeSchemaMarshaler)
	encoder.RegisterEncoder(DateTime{}, utils.EncodeSchemaMarshaler)
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *CustomRequest) QueryParams() *CustomRequestQueryParams {
	return r.queryParams
}

func (r CustomRequest) NewPathParams() *CustomRequestPathParams {
	return &CustomRequestPathParams{}
}

type CustomRequestPathParams struct {
}

func (p *CustomRequestPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *CustomRequest) PathParams() *CustomRequestPathParams {
	return r.pathParams
}

func (r *CustomRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *CustomRequest) SetMethod(method string) {
	r.method = method
}

func (r *CustomRequest) Method() string {
	return r.method
}

func (r CustomRequest) NewRequestBody() CustomRequestBody {
	return struct{}{}
}

type CustomRequestBody interface{}

func (r *CustomRequest) RequestBody() *CustomRequestBody {
	return &r.requestBody
}

func (r *CustomRequest) RequestBodyInterface() interface{} {
	return r.requestBody
}

func (r *CustomRequest) SetRequestBody(body CustomRequestBody) {
	r.requestBody = body
}

func (r *CustomRequest) NewResponseBody() *CustomResponseBody {
	var responseBody CustomResponseBody
	return &responseBody
}

type CustomResponseBody interface{}

func (r *CustomRequest) URL() (*url.URL, error) {
	u, err := r.client.GetEndpointURL("", r.PathParams())
	return &u, err
}

func (r *CustomRequest) Do() (CustomResponseBody, error) {
	// Create http request
	req, err := r.client.NewRequest(nil, r)
	if err != nil {
		return *r.NewResponseBody(), err
	}

	// Process query parameters
	err = utils.AddQueryParamsToRequest(r.QueryParams(), req, false)
	if err != nil {
		return *r.NewResponseBody(), err
	}

	responseBody := r.NewResponseBody()
	_, err = r.client.Do(req, responseBody)
	return *responseBody, err
}
