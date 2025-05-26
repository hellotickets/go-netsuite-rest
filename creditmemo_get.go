package netsuite

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/hellotickets/go-netsuite-rest/utils"
)

func (c *Client) NewCreditMemoGetRequest() CreditMemoGetRequest {
	r := CreditMemoGetRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type CreditMemoGetRequest struct {
	client      *Client
	queryParams *CreditMemoGetRequestQueryParams
	pathParams  *CreditMemoGetRequestPathParams
	method      string
	headers     http.Header
	requestBody CreditMemoGetRequestBody
}

func (r CreditMemoGetRequest) NewQueryParams() *CreditMemoGetRequestQueryParams {
	return &CreditMemoGetRequestQueryParams{}
}

type CreditMemoGetRequestQueryParams struct {
	Fields             Fields `schema:"fields,omitempty"`
	ExpandSubResources bool   `schema:"expandSubResources,omitempty"`
}

func (p CreditMemoGetRequestQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	encoder.RegisterEncoder(Date{}, utils.EncodeSchemaMarshaler)
	encoder.RegisterEncoder(DateTime{}, utils.EncodeSchemaMarshaler)
	encoder.RegisterEncoder(Fields{}, utils.EncodeSchemaMarshaler)
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *CreditMemoGetRequest) QueryParams() *CreditMemoGetRequestQueryParams {
	return r.queryParams
}

func (r CreditMemoGetRequest) NewPathParams() *CreditMemoGetRequestPathParams {
	return &CreditMemoGetRequestPathParams{}
}

type CreditMemoGetRequestPathParams struct {
	ID int `schema:"id"`
}

func (p *CreditMemoGetRequestPathParams) Params() map[string]string {
	return map[string]string{
		"id": strconv.Itoa(p.ID),
	}
}

func (r *CreditMemoGetRequest) PathParams() *CreditMemoGetRequestPathParams {
	return r.pathParams
}

func (r *CreditMemoGetRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *CreditMemoGetRequest) SetMethod(method string) {
	r.method = method
}

func (r *CreditMemoGetRequest) Method() string {
	return r.method
}

func (r CreditMemoGetRequest) NewRequestBody() CreditMemoGetRequestBody {
	return CreditMemoGetRequestBody{}
}

type CreditMemoGetRequestBody struct {
}

func (r *CreditMemoGetRequest) RequestBody() *CreditMemoGetRequestBody {
	return nil
}

func (r *CreditMemoGetRequest) RequestBodyInterface() interface{} {
	return nil
}

func (r *CreditMemoGetRequest) SetRequestBody(body CreditMemoGetRequestBody) {
	r.requestBody = body
}

func (r *CreditMemoGetRequest) NewResponseBody() *CreditMemoGetResponseBody {
	return &CreditMemoGetResponseBody{}
}

type CreditMemoGetResponseBody struct {
	Links Links `json:"links"`
	CreditMemo
}

func (r *CreditMemoGetRequest) URL() (*url.URL, error) {
	u, err := r.client.GetEndpointURL("/record/v1/creditMemo/{{.id}}", r.PathParams())
	return &u, err
}

func (r *CreditMemoGetRequest) Do() (CreditMemoGetResponseBody, error) {
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
