package holland

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"rentwatch/models"
)

type Holland struct {
	name   string
	SaasID string
}

func (h *Holland) Name() string {
	return h.name
}

type response struct {
	Units []struct {
		Beds  float64 `json:"beds,omitempty"`
		Baths float64 `json:"baths,omitempty"`
		Sqft  struct {
			Min float64 `json:"min,omitempty"`
			Max float64 `json:"max,omitempty"`
		} `json:"sqft"`
		Rent struct {
			Min float64 `json:"min,omitempty"`
			Max float64 `json:"max,omitempty"`
		} `json:"rent"`
	} `json:"units,omitempty"`
}

func (h *Holland) Units() ([]models.Unit, error) {
	u, err := url.Parse("https://www.hollandresidential.com/api/v1/content")

	query := u.Query()
	query.Add("saas_id", h.SaasID)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to send request: %v", err)
	}
	defer resp.Body.Close()

	result := &response{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("failed to json decode response: %v", err)
	}

	units := make([]models.Unit, 0)
	for _, unit := range result.Units {
		units = append(units, models.Unit{
			BathroomMin: unit.Baths,
			BathroomMax: unit.Baths,
			BedroomMin:  unit.Beds,
			BedroomMax:  unit.Beds,
			SqftMin:     unit.Sqft.Min,
			SqftMax:     unit.Sqft.Max,
			PriceMin:    unit.Rent.Min,
			PriceMax:    unit.Rent.Max,
		})
	}

	return units, nil
}

func NewHolland(name string, saasID string) *Holland {
	return &Holland{name: name, SaasID: saasID}
}

var Providers = []*Holland{
	NewHolland("Ivey on Boren", "Hiz4GtgjzZhE2rL4y"),
	NewHolland("Kiara", "v2MCbAhp2qPsB2GZg"),
	NewHolland("Dimension", "Gi8N4BuzsdrPYjmWL"),
	NewHolland("The Huxley", "mCpv7WScnT9XoYMLd"),
	NewHolland("One Lakefront", "aJRHh8cQ6cbHnq5JH"),
	NewHolland("JUXT", "RRfBiZh3PfLMjnwPA"),
}
