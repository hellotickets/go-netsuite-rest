package netsuite

import (
	"net/http"
	"net/url"

	"github.com/hellotickets/go-netsuite-rest/utils"
)

func (c *Client) NewUnitsTypeGetRequest() UnitsTypeGetRequest {
	r := UnitsTypeGetRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type UnitsTypeGetRequest struct {
	client      *Client
	queryParams *UnitsTypeGetRequestQueryParams
	pathParams  *UnitsTypeGetRequestPathParams
	method      string
	headers     http.Header
	requestBody UnitsTypeGetRequestBody
}

func (r UnitsTypeGetRequest) NewQueryParams() *UnitsTypeGetRequestQueryParams {
	return &UnitsTypeGetRequestQueryParams{}
}

type UnitsTypeGetRequestQueryParams struct {
	Q      string `schema:"q,omitempty"`
	Limit  int    `schema:"limit,omitempty"`
	Offset int    `schema:"offset,omitempty"`
}

func (p UnitsTypeGetRequestQueryParams) ToURLValues() (url.Values, error) {
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

func (r *UnitsTypeGetRequest) QueryParams() *UnitsTypeGetRequestQueryParams {
	return r.queryParams
}

func (r UnitsTypeGetRequest) NewPathParams() *UnitsTypeGetRequestPathParams {
	return &UnitsTypeGetRequestPathParams{}
}

type UnitsTypeGetRequestPathParams struct {
}

func (p *UnitsTypeGetRequestPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *UnitsTypeGetRequest) PathParams() *UnitsTypeGetRequestPathParams {
	return r.pathParams
}

func (r *UnitsTypeGetRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *UnitsTypeGetRequest) SetMethod(method string) {
	r.method = method
}

func (r *UnitsTypeGetRequest) Method() string {
	return r.method
}

func (r UnitsTypeGetRequest) NewRequestBody() UnitsTypeGetRequestBody {
	return UnitsTypeGetRequestBody{}
}

type UnitsTypeGetRequestBody struct {
}

func (r *UnitsTypeGetRequest) RequestBody() *UnitsTypeGetRequestBody {
	return nil
}

func (r *UnitsTypeGetRequest) RequestBodyInterface() interface{} {
	return nil
}

func (r *UnitsTypeGetRequest) SetRequestBody(body UnitsTypeGetRequestBody) {
	r.requestBody = body
}

func (r *UnitsTypeGetRequest) NewResponseBody() *UnitsTypeGetResponseBody {
	return &UnitsTypeGetResponseBody{}
}

type UnitsTypeGetResponseBody struct {
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

func (r *UnitsTypeGetRequest) URL() (*url.URL, error) {
	u, err := r.client.GetEndpointURL("/record/v1/unitsType", r.PathParams())
	return &u, err
}

func (r *UnitsTypeGetRequest) Do() (UnitsTypeGetResponseBody, error) {
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
