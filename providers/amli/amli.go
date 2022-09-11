package amli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rentwatch/models"
)

type Amli struct {
	name           string
	PropertyID     string
	AmliPropertyID string
}

func (a *Amli) Name() string {
	return a.name
}

type gql struct {
	Query         string         `json:"query,omitempty"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

type amliFloorplans struct {
	BathroomMin float64 `json:"bathroomMin,omitempty"`
	BathroomMax float64 `json:"bathroomMax,omitempty"`
	BedroomMin  float64 `json:"bedroomMin,omitempty"`
	BedroomMax  float64 `json:"bedroomMax,omitempty"`
	SqftMin     float64 `json:"sqftMin,omitempty"`
	SqftMax     float64 `json:"sqftMax,omitempty"`
	PriceMin    float64 `json:"priceMin,omitempty"`
	PriceMax    float64 `json:"priceMax,omitempty"`
}

type amliResponse struct {
	Data struct {
		PropertyFloorplansSummary []amliFloorplans `json:"propertyFloorplansSummary,omitempty"`
	} `json:"data"`
}

func (a *Amli) Units() ([]models.Unit, error) {
	query := gql{
		Query: `query Properties($amliPropertyId: ID!, $propertyId: ID!) {
    propertyFloorplansSummary(amliPropertyId: $amliPropertyId, propertyId: $propertyId) {
      bathroomMin
      bathroomMax
      bedroomMin
      bedroomMax
      priceMin
      priceMax
      sqftMin
      sqFtMax
    }
  }
`,
		OperationName: "Properties",
		Variables: map[string]any{
			"propertyId":     a.PropertyID,
			"amliPropertyId": a.AmliPropertyID,
		},
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(query)
	if err != nil {
		return nil, fmt.Errorf("unable to convert query to json: %v", err)
	}

	req, err := http.NewRequest("POST", "https://prodeastgraph.amli.com/graphql", buf)
	if err != nil {
		return nil, fmt.Errorf("unable to make request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("response error: %s", data)
	}

	result := &amliResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %v", err)
	}

	units := make([]models.Unit, 0)
	for _, f := range result.Data.PropertyFloorplansSummary {
		units = append(units, models.Unit{
			BathroomMin: f.BathroomMin,
			BathroomMax: f.BathroomMax,
			BedroomMin:  f.BedroomMin,
			BedroomMax:  f.BedroomMax,
			SqftMin:     f.SqftMin,
			SqftMax:     f.SqftMax,
			PriceMin:    f.PriceMin,
			PriceMax:    f.PriceMax,
		})
	}
	return units, nil
}

func NewAmli(name string, propertyID string, amliPropertyID string) *Amli {
	return &Amli{name: name, PropertyID: propertyID, AmliPropertyID: amliPropertyID}
}

var Providers = []*Amli{
	NewAmli("AMLI Arc", "XK-A2xAAAB8A_JrR", "89240"),
	NewAmli("AMLI Arc", "XK-A2xAAAB8A_JrR", "89240"),
	NewAmli("AMLI Wallingford", "XK-AAxAAACEA_JcG", "89178"),
	NewAmli("AMLI Mark24", "XHSPIxIAACIAbhlr", "88786"),
	NewAmli("AMLI 535", "XFJHfhMAACIANSgJ", "88146"),
	NewAmli("AMLI SLU", "XHSQBhIAAB8Abh1Y", "88848"),
	NewAmli("AMLI Bellevue Spring District", "XMxVmCwAADkA1DMw", "89407"),
	NewAmli("AMLI Bellevue Park", "XFJHVxMAACQANSdY", "85263"),
}
