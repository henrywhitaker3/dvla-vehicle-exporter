package dvla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientOptions struct {
	VesEndpoint string
	VesApiKey   string
}

func (c *ClientOptions) setDefaults() {
	if c.VesEndpoint == "" {
		c.VesEndpoint = "https://driver-vehicle-licensing.api.gov.uk"
	}
}

type Client struct {
	opts ClientOptions
}

func NewClient(opts ClientOptions) *Client {
	opts.setDefaults()
	return &Client{opts: opts}
}

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if time.Time(*d).IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, time.Time(*d).Format("2006-01-02"))), nil
}

type DateMonth time.Time

func (d *DateMonth) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		*d = DateMonth{}
		return nil
	}

	t, err := time.Parse("2006-01", s)
	if err != nil {
		return err
	}
	*d = DateMonth(t)
	return nil
}

func (d *DateMonth) MarshalJSON() ([]byte, error) {
	if time.Time(*d).IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, time.Time(*d).Format("2006-01"))), nil
}

type Vehicle struct {
	ArtEndDate               Date      `json:"artEndDate"`
	Co2Emissions             int       `json:"co2Emissions"`
	Colour                   string    `json:"colour"`
	EngineCapacity           int       `json:"engineCapacity"`
	FuelType                 string    `json:"fuelType"`
	Make                     string    `json:"make"`
	MarkedForExport          bool      `json:"markedForExport"`
	MonthOfFirstRegistration DateMonth `json:"monthOfFirstRegistration"`
	MotStatus                string    `json:"motStatus"`
	RegistrationNumber       string    `json:"registrationNumber"`
	RevenueWeight            int       `json:"revenueWeight"`
	TaxDueDate               Date      `json:"taxDueDate"`
	TaxStatus                string    `json:"taxStatus"`
	TypeApproval             string    `json:"typeApproval"`
	Wheelplan                string    `json:"wheelplan"`
	YearOfManufacture        int       `json:"yearOfManufacture"`
	EuroStatus               string    `json:"euroStatus"`
	RealDrivingEmissions     string    `json:"realDrivingEmissions"`
	DateOfLastV5CIssued      Date      `json:"dateOfLastV5CIssued"`
}

func (c *Client) GetVehicle(ctx context.Context, reg string) (*Vehicle, error) {
	body := struct {
		Reg string `json:"registrationNumber"`
	}{
		Reg: reg,
	}
	resp, err := c.vesRequest(ctx, http.MethodPost, "/vehicle-enquiry/v1/vehicles", body)
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

func (c *Client) vesRequest(ctx context.Context, method string, path string, body any) (*http.Response, error) {
	url, err := url.Parse(fmt.Sprintf("%s%s", c.opts.VesEndpoint, path))
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
	req.Header.Set("x-api-key", c.opts.VesApiKey)
	req = req.WithContext(ctx)

	return http.DefaultClient.Do(req)
}
