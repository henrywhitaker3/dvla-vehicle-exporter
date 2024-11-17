package dvla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ClientOptions struct {
	Endpoint string
	ApiKey   string
}

func (c *ClientOptions) setDefaults() {
	if c.Endpoint == "" {
		c.Endpoint = "https://driver-vehicle-licensing.api.gov.uk"
	}
}

type Client struct {
	opts ClientOptions
}

func NewClient(opts ClientOptions) *Client {
	opts.setDefaults()
	return &Client{opts: opts}
}

type Vehicle struct {
	ArtEndDate               time.Time `json:"artEndDate"`
	Co2Emissions             int       `json:"co2Emissions"`
	Colour                   string    `json:"colour"`
	EngineCapacity           int       `json:"engineCapacity"`
	FuelType                 string    `json:"fuelType"`
	Make                     string    `json:"make"`
	MarkedForExport          bool      `json:"markedForExport"`
	MonthOfFirstRegistration time.Time `json:"monthOfFirstRegistration"`
	MotStatus                string    `json:"motStatus"`
	RegistrationNumber       string    `json:"registrationNumber"`
	RevenueWeight            int       `json:"revenueWeight"`
	TaxDueDate               time.Time `json:"taxDueDate"`
	TaxStatus                string    `json:"taxStatus"`
	TypeApproval             string    `json:"typeApproval"`
	Wheelplan                string    `json:"wheelplan"`
	YearOfManufacture        int       `json:"yearOfManufacture"`
	EuroStatus               string    `json:"euroStatus"`
	RealDrivingEmissions     string    `json:"realDrivingEmissions"`
	DateOfLastV5CIssued      time.Time `json:"dateOfLastV5CIssued"`
}

func (c *Client) GetVehicle(ctx context.Context, reg string) (*Vehicle, error) {
	body := struct {
		Reg string `json:"registrationNumber"`
	}{
		Reg: reg,
	}
	resp, err := c.request(ctx, http.MethodPost, "/vehicle-enquiry/v1/vehicles", body)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle details: %w", err)
	}

	respBy, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle details response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get vehicle details status: %d body: %s", resp.StatusCode, string(respBy))
	}

	out := &Vehicle{}
	if err := json.Unmarshal(respBy, out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vehicle details response: %w", err)
	}
	return out, nil
}

func (c *Client) request(ctx context.Context, method string, path string, body any) (*http.Response, error) {
	url, err := url.Parse(fmt.Sprintf("%s%s", c.opts.Endpoint, path))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	var req *http.Request
	if body != nil {
		by, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body to json: %w", err)
		}
		req, err = http.NewRequest(method, url.String(), bytes.NewReader(by))
	} else {
		req, err = http.NewRequest(method, url.String(), nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("x-api-key", c.opts.ApiKey)
	req = req.WithContext(ctx)

	return http.DefaultClient.Do(req)
}
