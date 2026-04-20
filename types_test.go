package netsuite

import (
	"encoding/json"
	"testing"
)

// TestJSONNumber_AcceptsStringOrNumber documents the stdlib behavior the
// HT-14323 fix relies on: json.Number decodes both quoted and bare numeric
// JSON values. If a future stdlib change breaks this, the Invoice/CreditMemo
// sync will start failing again.
func TestJSONNumber_AcceptsStringOrNumber(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "quoted", in: `{"id":"777"}`, want: "777"},
		{name: "bare", in: `{"id":777}`, want: "777"},
		{name: "negative bare", in: `{"id":-1}`, want: "-1"},
		{name: "null", in: `{"id":null}`, want: ""},
		{name: "absent", in: `{}`, want: ""},
		{name: "empty string", in: `{"id":""}`, wantErr: true},
		{name: "non-numeric string", in: `{"id":"abc"}`, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s struct {
				ID json.Number `json:"id"`
			}
			err := json.Unmarshal([]byte(tc.in), &s)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got nil", tc.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("unmarshal %q: %v", tc.in, err)
			}
			if got := s.ID.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// TestInvoice_UnmarshalFixture reproduces the production decode failure:
// NetSuite sends item.items[0].item.id as a string, and itemSubtype/itemType
// as objects. Before the HT-14323 fix these fields failed to decode.
func TestInvoice_UnmarshalFixture(t *testing.T) {
	const body = `{
		"links": [],
		"id": "31457",
		"entity": {"id": "22784", "refName": "C-000096"},
		"item": {
			"links": [],
			"totalResults": 1,
			"items": [{
				"links": [],
				"account": {"id": "3031", "refName": "3031"},
				"amount": 1274.06,
				"description": "Test item",
				"item": {"id": "777", "refName": "EV-2120206"},
				"itemSubtype": {"id": "Sale", "refName": "Sale"},
				"itemType": {"id": "Service", "refName": "Service"},
				"quantity": 1.0,
				"taxAmount": 0.0,
				"taxDetailsReference": "31457_1"
			}]
		},
		"subsidiary": {"id": "1", "refName": "HELLO TICKET S.L."},
		"tranDate": "2025-01-23",
		"tranId": "F1INT000001"
	}`

	var inv Invoice
	if err := json.Unmarshal([]byte(body), &inv); err != nil {
		t.Fatalf("unmarshal invoice: %v", err)
	}

	if got := len(inv.Item.Items); got != 1 {
		t.Fatalf("expected 1 line item, got %d", got)
	}
	line := inv.Item.Items[0]
	if got := line.Item.ID.String(); got != "777" {
		t.Errorf("Item.ID: got %q, want %q", got, "777")
	}
	if got := line.ItemSubType.ID; got != "Sale" {
		t.Errorf("ItemSubType.ID: got %q, want %q", got, "Sale")
	}
	if got := line.ItemType.ID; got != "Service" {
		t.Errorf("ItemType.ID: got %q, want %q", got, "Service")
	}
}

// TestInvoice_UnmarshalNumericID exercises the defensive path where NetSuite
// might return a bare JSON number for nested ids.
func TestInvoice_UnmarshalNumericID(t *testing.T) {
	const body = `{
		"item": {
			"items": [{
				"item": {"id": 777, "refName": "x"}
			}]
		}
	}`
	var inv Invoice
	if err := json.Unmarshal([]byte(body), &inv); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := inv.Item.Items[0].Item.ID.String(); got != "777" {
		t.Errorf("Item.ID: got %q, want %q", got, "777")
	}
}

func TestCustomer_TaxRegistrationIDString(t *testing.T) {
	const body = `{
		"id": "22784",
		"taxRegistration": {
			"links": [],
			"totalResults": 1,
			"items": [{
				"links": [],
				"id": "123",
				"nexus": {"id": "1", "refName": "Spain"},
				"nexusCountry": {"id": "ES", "refName": "Spain"},
				"taxRegistrationNumber": "B87643680"
			}]
		}
	}`

	var c Customer
	if err := json.Unmarshal([]byte(body), &c); err != nil {
		t.Fatalf("unmarshal customer: %v", err)
	}
	if got := len(c.TaxRegistration.Items); got != 1 {
		t.Fatalf("expected 1 tax registration item, got %d", got)
	}
	if got := c.TaxRegistration.Items[0].ID.String(); got != "123" {
		t.Errorf("TaxRegistration.Items[0].ID: got %q, want %q", got, "123")
	}
	if got := c.TaxRegistration.Items[0].TaxRegistrationNumber; got != "B87643680" {
		t.Errorf("TaxRegistrationNumber: got %q", got)
	}
}

// TestErrorResponse_ObjectType verifies ErrorResponse tolerates NetSuite
// returning `type` as an object (not just the documented URL string).
func TestErrorResponse_ObjectType(t *testing.T) {
	const body = `{
		"type": {"code": "INSUFFICIENT_PERMISSION"},
		"title": "Forbidden",
		"status": 403,
		"o:errorDetails": [{"detail": "d", "o:errorCode": "E"}]
	}`

	var er ErrorResponse
	if err := json.Unmarshal([]byte(body), &er); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if got := er.Error(); got != "E: d" {
		t.Errorf("Error(): got %q, want %q", got, "E: d")
	}
}

func TestErrorResponse_StringType(t *testing.T) {
	const body = `{
		"type": "https://example.com/problems/forbidden",
		"title": "Forbidden",
		"status": 403,
		"o:errorDetails": [{"detail": "d", "o:errorCode": "E"}]
	}`

	var er ErrorResponse
	if err := json.Unmarshal([]byte(body), &er); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if got := er.Error(); got != "E: d" {
		t.Errorf("Error(): got %q, want %q", got, "E: d")
	}
}
