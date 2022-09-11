package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"rentwatch/models"
	"rentwatch/providers/amli"
	"rentwatch/providers/holland"
	"rentwatch/providers/sightmap"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type statistics struct {
	Date                  string
	TotalProviders        int
	TotalProvidersSuccess int
	TotalProvidersError   int
	TotalUnits            int
}

func main() {
	if err := applyMigrations("sqlite3://rentwatch.db"); err != nil {
		log.Fatalf("unable to apply migrations: %v", err)
	}

	now := time.Now().UTC()

	db, err := sql.Open("sqlite3", "./rentwatch.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Slices can't splat into interfaces :/
	providers := loadProviders()

	stats := statistics{
		Date:           now.Format(time.RFC3339),
		TotalProviders: len(providers),
	}

	for _, provider := range providers {
		builder := sq.Insert("units").Columns(
			"crawl_date",
			"name",
			"bed_min",
			"bed_max",
			"bath_min",
			"bath_max",
			"sqft_min",
			"sqft_max",
			"price_min",
			"price_max",
		)
		units, err := provider.Units()
		if err != nil {
			log.Printf("failed %s: %v", provider.Name(), err)
			stats.TotalProvidersError += 1
			continue
		}
		for _, unit := range units {
			builder = builder.Values(
				now.Format(time.RFC3339),
				provider.Name(),
				unit.BedroomMin,
				unit.BedroomMax,
				unit.BathroomMin,
				unit.BathroomMax,
				unit.SqftMin,
				unit.SqftMax,
				unit.PriceMin,
				unit.PriceMax,
			)
		}
		query, args, err := builder.ToSql()
		if err != nil {
			log.Printf("failed to build query for %s: %v", provider.Name(), err)
			stats.TotalProvidersError += 1
			continue
		}
		_, err = db.Exec(query, args...)
		if err != nil {
			log.Printf("failed to insert units for %s: %v", provider.Name(), err)
			stats.TotalProvidersError += 1
			continue
		}
		stats.TotalUnits += len(units)
		stats.TotalProvidersSuccess += 1
	}
	if err := json.NewEncoder(os.Stderr).Encode(stats); err != nil {
		log.Fatalf("unable to display stats: %v", err)
	}
}

func loadProviders() []models.Provider {
	providers := make([]models.Provider, 0)
	// Load amli providers.
	for _, provider := range amli.Providers {
		providers = append(providers, provider)
	}

	// Load holland providers.
	for _, provider := range holland.Providers {
		providers = append(providers, provider)
	}

	// Load sightmap providers.
	for _, provider := range sightmap.Providers {
		providers = append(providers, provider)
	}
	return providers
}

func applyMigrations(dsn string) error {
	m, err := migrate.New("file://./migrations", dsn)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
