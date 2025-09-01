package nominatim

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type api interface {
	ExecuteNominatimRequest(r *http.Request) ([]NominatimCoreResult, error)
	RequestWithSearchString(ctx context.Context, url, searchString, limit, language string) ([]NominatimCoreResult, error)
	RequestWithParameters(ctx context.Context, url string, parameters NominatimDetailRequest, limit string, language string) ([]NominatimCoreResult, error)
}

type apiNominatim struct {
	client    http.Client
	userAgent string
}

func (napi apiNominatim) ExecuteNominatimRequest(r *http.Request) ([]NominatimCoreResult, error) {
	r.Header.Add("User-Agent", napi.userAgent)
	res, err := napi.client.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code not 200, instead: %d", res.StatusCode)
	}

	if res.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	var result []NominatimResult
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		// unable to unmarshal json
		return nil, err
	}

	var tmp []NominatimCoreResult
	for _, i := range result {
		tmp = append(tmp, i.Address)
	}

	return tmp, nil
}

func (napi apiNominatim) RequestWithSearchString(ctx context.Context, url, searchString, limit, language string) ([]NominatimCoreResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if len(limit) > 0 {
		q.Add("limit", limit)
	}
	q.Add("accept-language", language)
	q.Add("addressdetails", "1")
	q.Add("format", "json")
	q.Add("q", searchString)
	req.URL.RawQuery = q.Encode()
	return napi.ExecuteNominatimRequest(req)
}

type NominatimDetailRequest struct {
	Street     string
	PostalCode string
	City       string
}

func (napi apiNominatim) RequestWithParameters(ctx context.Context, url string, parameters NominatimDetailRequest, limit string, language string) ([]NominatimCoreResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if len(limit) > 0 {
		q.Add("limit", limit)
	}
	q.Add("accept-language", language)
	q.Add("addressdetails", "1")
	q.Add("format", "json")
	if len(parameters.Street) > 0 {
		q.Add("street", parameters.Street)
	}
	if len(parameters.City) > 0 {
		q.Add("city", parameters.City)
	}
	if len(parameters.PostalCode) > 0 {
		q.Add("postalcode", parameters.PostalCode)
	}

	req.URL.RawQuery = q.Encode()

	return napi.ExecuteNominatimRequest(req)
}
