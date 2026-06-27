package mijnhost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DNSControl/dnscontrol/v4/providers/mijnhost/mijnhostapi"
)

const (
	BASE_ENDPOINT = "https://mijn.host/api/v2"
	USER_AGENT    = "DNSControl (providers/mijnhost v0.1)"
)

type mijnhostAPIClient struct {
	httpClient *http.Client

	apikey string
}

func newAPIClient(apikey string, timeout time.Duration) *mijnhostAPIClient {
	return &mijnhostAPIClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		apikey: apikey,
	}
}

func (ac *mijnhostAPIClient) GetDomain(domain string) (mijnhostapi.GetDomainResponseData, error) {
	var apiResponse mijnhostapi.GetDomainResponse

	req, err := ac.newRequest(http.MethodGet, "/domains/"+domain, nil)
	if err != nil {
		return apiResponse.Data, fmt.Errorf("GetDomain request creation failed: %w", err)
	}

	resp, err := ac.performRequest(req)
	if err != nil {
		return apiResponse.Data, fmt.Errorf("GetDomain performing request failed: %w", err)
	}
	defer resp.Body.Close()

	err = ac.convertResponse(resp, &apiResponse)
	if err != nil {
		return apiResponse.Data, fmt.Errorf("GetDomain converting response body failed: %w", err)
	}

	return apiResponse.Data, nil
}

func (ac *mijnhostAPIClient) GetDNSRecordsForDomain(domain string) (mijnhostapi.GetDNSRecordsForDomainResponseData, error) {
	var apiResponse mijnhostapi.GetDNSRecordsForDomainResponse

	req, err := ac.newRequest(http.MethodGet, "/domains/"+domain+"/dns", nil)
	if err != nil {
		return apiResponse.Data, fmt.Errorf("GetDNSRecordsForDomain request creation failed: %w", err)
	}

	resp, err := ac.performRequest(req)
	if err != nil {
		return apiResponse.Data, fmt.Errorf("GetDNSRecordsForDomain performing request failed: %w", err)
	}
	defer resp.Body.Close()

	err = ac.convertResponse(resp, &apiResponse)
	if err != nil {
		return apiResponse.Data, fmt.Errorf("GetDNSRecordsForDomain converting response body failed: %w", err)
	}

	return apiResponse.Data, nil
}

func (ac *mijnhostAPIClient) UpdateDNSRecordForDomain(domain string, rec mijnhostapi.Record) (mijnhostapi.UpdateDNSRecordForDomainResponse, error) {
	var apiResponse mijnhostapi.UpdateDNSRecordForDomainResponse

	body, err := json.Marshal(mijnhostapi.PatchRecord{rec})
	if err != nil {
		return apiResponse, fmt.Errorf("UpdateDNSRecordForDomain marshaling request failed: %w", err)
	}

	req, err := ac.newRequest(http.MethodPatch, "/domains/"+domain+"/dns", bytes.NewBuffer(body))
	if err != nil {
		return apiResponse, fmt.Errorf("UpdateDNSRecordForDomain request creation failed: %w", err)
	}

	resp, err := ac.performRequest(req)
	if err != nil {
		return apiResponse, fmt.Errorf("UpdateDNSRecordForDomain performing request failed: %w", err)
	}
	defer resp.Body.Close()

	err = ac.convertResponse(resp, &apiResponse)
	if err != nil {
		return apiResponse, fmt.Errorf("UpdateDNSRecordForDomain converting response body failed: %w", err)
	}

	return apiResponse, nil
}

func (ac *mijnhostAPIClient) DeleteDNSRecordForDomain(domain string, rec mijnhostapi.Record) (mijnhostapi.DeleteDNSRecordForDomainResponse, error) {
	var apiResponse mijnhostapi.DeleteDNSRecordForDomainResponse

	body, err := json.Marshal(mijnhostapi.DeleteRecord{rec})
	if err != nil {
		return apiResponse, fmt.Errorf("DeleteDNSRecordForDomain marshaling request failed: %w", err)
	}

	req, err := ac.newRequest(http.MethodDelete, "/domains/"+domain+"/dns", bytes.NewBuffer(body))
	if err != nil {
		return apiResponse, fmt.Errorf("DeleteDNSRecordForDomain request creation failed: %w", err)
	}

	resp, err := ac.performRequest(req)
	if err != nil {
		return apiResponse, fmt.Errorf("DeleteDNSRecordForDomain performing request failed: %w", err)
	}
	defer resp.Body.Close()

	err = ac.convertResponse(resp, &apiResponse)
	if err != nil {
		return apiResponse, fmt.Errorf("DeleteDNSRecordForDomain converting response body failed: %w", err)
	}

	return apiResponse, nil
}

func (ac *mijnhostAPIClient) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, BASE_ENDPOINT+"/"+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("API-Key", ac.apikey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// performRequest executes the supplied request. The caller must close the response body on nil error.
func (ac *mijnhostAPIClient) performRequest(req *http.Request) (*http.Response, error) {
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	// TODO: handle http status codes

	return resp, nil
}

// convertResponse checks the supported type for target, and unmarshals resp.Body into it.
func (ac *mijnhostAPIClient) convertResponse(resp *http.Response, target any) error {
	switch target.(type) {
	case *mijnhostapi.GetDomainResponse:
		// OK
	case *mijnhostapi.GetDNSRecordsForDomainResponse:
		// OK
	case *mijnhostapi.UpdateDNSRecordForDomainResponse:
		// OK
	case *mijnhostapi.DeleteDNSRecordForDomainResponse:
		// OK
	default:
		panic("unsupported API response type")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("convertResponse ioutil.ReadAll failed: %w", err)
	}

	err = json.Unmarshal(body, &target)
	if err != nil {
		return fmt.Errorf("convertResponse json.Unmarshal failed: %w", err)
	}

	return nil
}
