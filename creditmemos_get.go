package netsuite

import (
	"net/http"
	"net/url"

	"github.com/hellotickets/go-netsuite-rest/utils"
)

func (c *Client) NewCreditMemosGetRequest() CreditMemosGetRequest {
	r := CreditMemosGetRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type CreditMemosGetRequest struct {
	client      *Client
	queryParams *CreditMemosGetRequestQueryParams
	pathParams  *CreditMemosGetRequestPathParams
	method      string
	headers     http.Header
	requestBody CreditMemosGetRequestBody
}

func (r CreditMemosGetRequest) NewQueryParams() *CreditMemosGetRequestQueryParams {
	return &CreditMemosGetRequestQueryParams{}
}

type CreditMemosGetRequestQueryParams struct {
	Q      string `schema:"q,omitempty"`
	Limit  int    `schema:"limit,omitempty"`
	Offset int    `schema:"offset,omitempty"`
}

func (p CreditMemosGetRequestQueryParams) ToURLValues() (url.Values, error) {
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

func (r *CreditMemosGetRequest) QueryParams() *CreditMemosGetRequestQueryParams {
	return r.queryParams
}

func (r CreditMemosGetRequest) NewPathParams() *CreditMemosGetRequestPathParams {
	return &CreditMemosGetRequestPathParams{}
}

type CreditMemosGetRequestPathParams struct {
}

func (p *CreditMemosGetRequestPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *CreditMemosGetRequest) PathParams() *CreditMemosGetRequestPathParams {
	return r.pathParams
}

func (r *CreditMemosGetRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *CreditMemosGetRequest) SetMethod(method string) {
	r.method = method
}

func (r *CreditMemosGetRequest) Method() string {
	return r.method
}

func (r CreditMemosGetRequest) NewRequestBody() CreditMemosGetRequestBody {
	return CreditMemosGetRequestBody{}
}

type CreditMemosGetRequestBody struct {
}

func (r *CreditMemosGetRequest) RequestBody() *CreditMemosGetRequestBody {
	return nil
}

func (r *CreditMemosGetRequest) RequestBodyInterface() interface{} {
	return nil
}

func (r *CreditMemosGetRequest) SetRequestBody(body CreditMemosGetRequestBody) {
	r.requestBody = body
}

func (r *CreditMemosGetRequest) NewResponseBody() *CreditMemosGetResponseBody {
	return &CreditMemosGetResponseBody{}
}

type CreditMemosGetResponseBody struct {
	Links []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
	Count   int  `json:"count"`
	HasMore bool `json:"hasMore"`
	Items   []struct {
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
		ID string `json:"id"`
	} `json:"items"`
	Offset       int `json:"offset"`
	TotalResults int `json:"totalResults"`
}

func (r *CreditMemosGetRequest) URL() (*url.URL, error) {
	u, err := r.client.GetEndpointURL("/record/v1/creditMemo", r.PathParams())
	return &u, err
}

func (r *CreditMemosGetRequest) Do() (CreditMemosGetResponseBody, error) {
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
