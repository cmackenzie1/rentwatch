package sightmap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rentwatch/models"
)

type Sightmap struct {
	name string
	URL  string
}

func (s *Sightmap) Name() string {
	return s.name
}

type floorPlan struct {
	ID            string  `json:"id,omitempty"`
	Name          string  `json:"name,omitempty"`
	BedroomCount  float64 `json:"bedroom_count,omitempty"`
	BathroomCount float64 `json:"bathroom_count,omitempty"`
}

type unit struct {
	FloorPlanID string  `json:"floor_plan_id,omitempty"`
	Area        float64 `json:"area,omitempty"`
	Price       float64 `json:"price,omitempty"`
}

type response struct {
	Data struct {
		FloorPlans []floorPlan `json:"floor_plans,omitempty"`
		Units      []unit      `json:"units,omitempty"`
	} `json:"data,omitempty"`
}

func (s *Sightmap) Units() ([]models.Unit, error) {
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	results := &response{}
	err = json.NewDecoder(resp.Body).Decode(results)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	floorPlans := make(map[string]floorPlan)
	for _, plan := range results.Data.FloorPlans {
		floorPlans[plan.ID] = plan
	}

	units := make([]models.Unit, 0)
	for _, u := range results.Data.Units {
		units = append(units, models.Unit{
			BathroomMin: floorPlans[u.FloorPlanID].BathroomCount,
			BathroomMax: floorPlans[u.FloorPlanID].BathroomCount,
			BedroomMin:  floorPlans[u.FloorPlanID].BedroomCount,
			BedroomMax:  floorPlans[u.FloorPlanID].BedroomCount,
			SqftMin:     u.Area,
			SqftMax:     u.Area,
			PriceMin:    u.Price,
			PriceMax:    u.Price,
		})
	}

	return units, nil
}

func NewSightmap(name string, URL string) *Sightmap {
	return &Sightmap{name: name, URL: URL}
}

var Providers = []*Sightmap{
	NewSightmap("REN", "https://sightmap.com/app/api/v1/l8xvrjnmpjk/sightmaps/13986"),
	NewSightmap("McKenzie", "https://sightmap.com/app/api/v1/6m9pzykzvk1/sightmaps/1106"),
	NewSightmap("888 Bellevue", "https://sightmap.com/app/api/v1/dzlporo4pg4/sightmaps/5526"),
	NewSightmap("Broadstone Sky", "https://sightmap.com/app/api/v1/9zw467jlp87/sightmaps/23139"),
}
