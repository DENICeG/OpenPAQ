package nominatim

import (
	"reflect"
	"testing"
)

func TestNominatimCoreResult_parse(t *testing.T) {
	type fields struct {
		State            string
		StateDistrict    string
		Municipality     string
		City             string
		Town             string
		Village          string
		CityDistrict     string
		District         string
		Borough          string
		Suburb           string
		Subdivision      string
		Hamlet           string
		Croft            string
		IsolatedDwelling string
		Neighbourhood    string
		Allotments       string
		Quarter          string
		Residential      string
		Farm             string
		Farmyard         string
		Industrial       string
		Commercial       string
		Retail           string
		Road             string
		CityBlock        string
		PostCode         string
		CountryCode      string
	}

	raw := fields{
		State:            "Germany",
		StateDistrict:    "NRW",
		Municipality:     "Dortmund",
		City:             "Dortmund",
		Town:             "Dortmund Nord",
		Village:          "Dortmund Nordkaff",
		CityDistrict:     "Dortmund N1",
		District:         "N1",
		Borough:          "Dborough",
		Suburb:           "Suburb",
		Subdivision:      "Subdivision",
		Hamlet:           "Hamlet",
		Croft:            "Croft",
		IsolatedDwelling: "Isolated",
		Neighbourhood:    "Neighbourhood",
		Allotments:       "Allot",
		Quarter:          "Quarter",
		Residential:      "Resident",
		Farm:             "Farm 1",
		Farmyard:         "Yard",
		Industrial:       "Industrial Area",
		Commercial:       "Commerc",
		Retail:           "Retail",
		Road:             "Road",
		CityBlock:        "CB",
		PostCode:         "DE-12345",
		CountryCode:      "DE",
	}
	expectedResult := ParsedResult{
		Street:      []string{"hamlet", "road", "cb"},
		PostalCode:  "de-12345",
		City:        []string{"germany", "nrw", "dortmund", "dortmund", "dortmund nord", "dortmund nordkaff", "dortmund n1", "n1", "dborough", "suburb", "subdivision", "hamlet", "croft", "isolated", "neighbourhood", "allot", "quarter", "resident", "farm 1", "yard", "industrial area", "commerc", "retail"},
		CountryCode: "de",
	}

	tests := []struct {
		name   string
		fields fields
		want   ParsedResult
	}{
		{name: "parse-all", fields: raw, want: expectedResult},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nr := &NominatimCoreResult{
				State:            tt.fields.State,
				StateDistrict:    tt.fields.StateDistrict,
				Municipality:     tt.fields.Municipality,
				City:             tt.fields.City,
				Town:             tt.fields.Town,
				Village:          tt.fields.Village,
				CityDistrict:     tt.fields.CityDistrict,
				District:         tt.fields.District,
				Borough:          tt.fields.Borough,
				Suburb:           tt.fields.Suburb,
				Subdivision:      tt.fields.Subdivision,
				Hamlet:           tt.fields.Hamlet,
				Croft:            tt.fields.Croft,
				IsolatedDwelling: tt.fields.IsolatedDwelling,
				Neighbourhood:    tt.fields.Neighbourhood,
				Allotments:       tt.fields.Allotments,
				Quarter:          tt.fields.Quarter,
				Residential:      tt.fields.Residential,
				Farm:             tt.fields.Farm,
				Farmyard:         tt.fields.Farmyard,
				Industrial:       tt.fields.Industrial,
				Commercial:       tt.fields.Commercial,
				Retail:           tt.fields.Retail,
				Road:             tt.fields.Road,
				CityBlock:        tt.fields.CityBlock,
				PostCode:         tt.fields.PostCode,
				CountryCode:      tt.fields.CountryCode,
			}
			if got := nr.parse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
