// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package de

import (
	"context"
	"fmt"
	"openPAQ/internal/listmatcher/types"
	"openPAQ/internal/normalization"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Database interface {
	GetPostalCodeStreet(normalizer normalization.Normalize) (map[string]PostalCodeStreetItems, error)
	GetCityPostalCode(normalizer normalization.Normalize) (map[string]CityPostalCodeItems, error)
}

type ClickhouseDatabase struct {
	db       clickhouse.Conn
	database string
	table    string
}

type LookupList struct {
	Name       string `ch:"strasse_name"`
	PostalCode string `ch:"postleitzahl"`
	Gemeinde   string `ch:"gemeinde_name"`
	Ort        string `ch:"ort_name"`
}

func NewDatabase(config types.DatabaseConfig) *ClickhouseDatabase {
	return &ClickhouseDatabase{
		db:       initClickhouseConfig(config),
		database: config.DataBase,
		table:    config.Table,
	}
}

func initClickhouseConfig(config types.DatabaseConfig) (db clickhouse.Conn) {
	db, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", config.DbHost, config.DbPort)},
		Auth: clickhouse.Auth{
			Database: config.DataBase,
			Username: config.DbUserName,
			Password: config.DbUserPassword,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("error connecting to database: %s", err))
	}
	return db
}

func (db *ClickhouseDatabase) GetPostalCodeStreet(normalizer normalization.Normalize) (map[string]PostalCodeStreetItems, error) {

	result := make(map[string]PostalCodeStreetItems)

	type LookupList struct {
		PostalCode string `ch:"postleitzahl"`
		Street     string `ch:"strasse_name"`
	}

	var dbResponse []LookupList

	query := `
		SELECT DISTINCT postleitzahl,strasse_name
		FROM {database:Identifier}.{table:Identifier} 
		ORDER BY postleitzahl
	`

	ctx := clickhouse.Context(context.Background(), clickhouse.WithParameters(clickhouse.Parameters{
		"database": db.database,
		"table":    db.table,
	}))

	err := db.db.Select(ctx, &dbResponse, query)

	if err != nil {
		return nil, err
	}

	for _, entry := range dbResponse {

		normalizedPostalCode, err := normalizer.PostalCode(entry.PostalCode)
		if err != nil {
			continue
		}

		normalizedStreets, err := normalizer.Street(entry.Street)
		if err != nil {
			continue
		}

		postalCodeStreetItem, ok := result[entry.PostalCode]

		var streets []NormalizeEnclosure

		for _, normalizedStreet := range normalizedStreets {
			streets = append(streets, NormalizeEnclosure{
				Raw:        entry.Street,
				Normalized: normalizedStreet,
			})
		}

		if ok {
			streets = append(streets, postalCodeStreetItem.Streets...)
		}

		result[normalizedPostalCode] = PostalCodeStreetItems{
			PostalCode: entry.PostalCode,
			Streets:    streets,
		}

	}

	return result, nil
}

func (db *ClickhouseDatabase) GetCityPostalCode(normalizer normalization.Normalize) (map[string]CityPostalCodeItems, error) {

	result := make(map[string]CityPostalCodeItems)

	type LookupList struct {
		Region     string `ch:"gemeinde_name"`
		City       string `ch:"ort_name"`
		PostalCode string `ch:"postleitzahl"`
	}

	var dbResponse []LookupList

	query := `
		SELECT DISTINCT gemeinde_name,ort_name,postleitzahl
		FROM {database:Identifier}.{table:Identifier} 
		ORDER BY gemeinde_name,ort_name
	`

	ctx := clickhouse.Context(context.Background(), clickhouse.WithParameters(clickhouse.Parameters{
		"database": db.database,
		"table":    db.table,
	}))

	err := db.db.Select(ctx, &dbResponse, query)

	if err != nil {
		return nil, err
	}

	for _, entry := range dbResponse {

		normalizedPostalCode, err := normalizer.PostalCode(entry.PostalCode)
		if err != nil || normalizedPostalCode == "" {
			continue
		}

		postalCodeItem := NormalizeEnclosure{
			Raw:        entry.PostalCode,
			Normalized: normalizedPostalCode,
		}

		normalizedCity, err := normalizer.City(entry.City)
		if err != nil {
			continue
		}

		cityPostalCodeItem, ok := result[normalizedCity]

		if !ok {
			result[normalizedCity] = CityPostalCodeItems{
				City: entry.City,
				PostalCodes: []NormalizeEnclosure{
					postalCodeItem,
				},
			}
		} else {
			temp := cityPostalCodeItem.PostalCodes
			temp = append(temp, postalCodeItem)

			result[normalizedCity] = CityPostalCodeItems{
				City:        result[normalizedCity].City,
				PostalCodes: temp,
			}

		}

		normalizedRegion, err := normalizer.City(entry.Region)
		if err != nil {
			continue
		}

		cityPostalCodeItem, ok = result[normalizedRegion]

		if !ok {
			result[normalizedRegion] = CityPostalCodeItems{
				City: entry.Region,
				PostalCodes: []NormalizeEnclosure{
					postalCodeItem,
				},
			}
		} else {
			temp := cityPostalCodeItem.PostalCodes
			found := false
			for _, v := range temp {
				if v.Normalized == normalizedPostalCode {
					found = true
				}
			}

			if !found {

				temp = append(temp, postalCodeItem)

				result[normalizedRegion] = CityPostalCodeItems{
					City:        result[normalizedRegion].City,
					PostalCodes: temp,
				}
			}

		}

	}

	return result, nil
}
