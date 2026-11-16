// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"fmt"
	"regexp"
	"strings"
)

type ES struct {
	rePostalCode        *regexp.Regexp
	reCrap              *regexp.Regexp
	reStreetparts       *regexp.Regexp
	reStreetShortName   *regexp.Regexp
	reStreetHouseNumber *regexp.Regexp
	reCity              *regexp.Regexp
}

func newES() (*ES, error) {
	rePostalCode, errPostalCode := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;a-zA-Z- ]")
	if errPostalCode != nil {
		return nil, errPostalCode
	}

	reCrap, errCrap := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9]")
	if errCrap != nil {
		return nil, errCrap
	}
	reCity, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;0-9\n\r]")
	if err != nil {
		return nil, err
	}
	//street etc. just in spanish
	reStreetparts, errStreetparts := regexp.Compile("\\b(autopista|autoviac|avenida|bulevar|calle( de)?( la)?|calle peatonal|carrer( de)?( la)?|callejon|camino|canada real|carretera|carretera de circunvalacion|carril|ciclovia|corredera|costanilla|parque|pasadizo elevado|pasaje|paseo maritimo|plaza|pretil|puente|ronda|sendero|travesia|tunel|via pecuaria|via rapida|via verde|urbanizacion)\\b")
	if errStreetparts != nil {
		return nil, errStreetparts
	}

	reStreetShortName, errStreetShortName := regexp.Compile("\\b(c[/|.]|av[/|.]|avda[/|.]|pl[/|.]|s/n|sin numero|n[º|.]|numero)")

	if errStreetShortName != nil {
		return nil, errStreetparts
	}

	reStreetHouseNumber, errreStreetHouseNumber := regexp.Compile("\\b(\\d[º|-]\\w)\\b")

	if errreStreetHouseNumber != nil {
		return nil, errStreetparts
	}
	return &ES{
		rePostalCode:        rePostalCode,
		reCrap:              reCrap,
		reStreetparts:       reStreetparts,
		reStreetShortName:   reStreetShortName,
		reStreetHouseNumber: reStreetHouseNumber,
		reCity:              reCity,
	}, nil
}
func (es *ES) ReplaceSpanishLetters(s string) string {
	s = strings.ReplaceAll(s, "í", "i")
	s = strings.ReplaceAll(s, "ç", "c")
	s = strings.ReplaceAll(s, "ó", "o")
	s = strings.ReplaceAll(s, "ñ", "n")
	s = strings.ReplaceAll(s, "ú", "u")
	s = strings.ReplaceAll(s, "á", "a")
	return s
}

func (es *ES) GetCountryCode() string {
	return "es"
}

func (es *ES) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = es.reCity.ReplaceAllString(s, "")
	s = es.ReplaceSpanishLetters(s)

	return s, nil
}

func (es *ES) PostalCode(s string) (string, error) {
	s = es.rePostalCode.ReplaceAllString(s, "")

	if len(s) != 5 {
		if len(s) > 5 {
			return s[:5], fmt.Errorf("not valid postalcode")
		}
		return s, fmt.Errorf("not valid postalcode")
	}

	return s, nil
}

func (es *ES) Street(s string) ([]string, error) {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "\n", "|")
	s = strings.ReplaceAll(s, ",", "|")
	s = es.ReplaceSpanishLetters(s)
	s = es.reStreetShortName.ReplaceAllString(s, " ")
	s = es.reStreetparts.ReplaceAllString(s, " ")
	s = es.reStreetHouseNumber.ReplaceAllString(s, " ")
	s = es.reCrap.ReplaceAllString(s, " ")

	addressSlice := strings.Split(s, "|")
	var cleanAddressSlice []string
	for _, v := range addressSlice {
		var cleanedParts []string
		addressPartSlices := strings.Split(v, " ")
		for _, p := range addressPartSlices {
			if len(p) > 1 {
				cleanedParts = append(cleanedParts, p)
			}
		}
		cleanPart := strings.Join(cleanedParts, " ")
		if len(cleanPart) > 1 {
			cleanAddressSlice = append(cleanAddressSlice, cleanPart)
		}
	}

	return cleanAddressSlice, nil
}
