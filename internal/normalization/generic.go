// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"regexp"
	"strings"
)

type Generic struct {
	reStreet     *regexp.Regexp
	rePostalCode *regexp.Regexp
	reCity       *regexp.Regexp
}

func NewGeneric() (*Generic, error) {
	reStreet, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9-\n\r]")
	if err != nil {
		return nil, err
	}

	rePostalCode, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;\n\r]")
	if err != nil {
		return nil, err
	}

	reCity, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;0-9\n\r]")
	if err != nil {
		return nil, err
	}
	return &Generic{
		reStreet:     reStreet,
		rePostalCode: rePostalCode,
		reCity:       reCity,
	}, nil
}

func (g *Generic) GetCountryCode() string {
	return "generic"
}

func (g *Generic) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = g.reCity.ReplaceAllString(s, "")
	return strings.TrimSpace(s), nil
}

func (g *Generic) PostalCode(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	return g.rePostalCode.ReplaceAllString(s, ""), nil
}

func (g *Generic) Street(s string) ([]string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ",", "|")
	s = strings.ReplaceAll(s, "\n", "|")
	s = g.reStreet.ReplaceAllString(s, "")

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
