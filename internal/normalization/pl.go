// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type PL struct {
	reStreet              *regexp.Regexp
	reStreetDeleteWords   *regexp.Regexp
	rePostalCode          *regexp.Regexp
	rePostalCodeIsNumeric *regexp.Regexp
	reCity                *regexp.Regexp
	reWordEndWithDot      *regexp.Regexp
}

func NewPl() (*PL, error) {
	reStreet, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9-\n\r]")
	if err != nil {
		return nil, err
	}

	reStreetDeleteWords, err := regexp.Compile("\\b(ul|al|pl|os)\\b")
	if err != nil {
		return nil, err
	}

	reWordEndWithDot, err := regexp.Compile("\\b[\\w]*\\.")
	if err != nil {
		return nil, err
	}

	rePostalCode, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;\n\rA-Za-z]")
	if err != nil {
		return nil, err
	}

	rePostalCodeIsNumeric, err := regexp.Compile("^[\\d]{5}$")
	if err != nil {
		return nil, err
	}

	reCity, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€'¿_\",:;\n\r0-9]")
	if err != nil {
		return nil, err
	}
	return &PL{
		reStreet:              reStreet,
		reStreetDeleteWords:   reStreetDeleteWords,
		rePostalCode:          rePostalCode,
		rePostalCodeIsNumeric: rePostalCodeIsNumeric,
		reCity:                reCity,
		reWordEndWithDot:      reWordEndWithDot,
	}, nil
}

func (p *PL) GetCountryCode() string {
	return "pl"
}

func replaceSpecialChars(s string) string {
	s = strings.ReplaceAll(s, "ę", "e")
	s = strings.ReplaceAll(s, "ó", "o")
	s = strings.ReplaceAll(s, "ą", "a")
	s = strings.ReplaceAll(s, "ł", "l")
	s = strings.ReplaceAll(s, "ż", "z")
	s = strings.ReplaceAll(s, "ź", "z")
	s = strings.ReplaceAll(s, "ś", "s")
	s = strings.ReplaceAll(s, "ć", "c")
	return s
}

func (p *PL) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = p.reCity.ReplaceAllString(s, "")
	s = replaceSpecialChars(s)
	s = strings.TrimSpace(s)

	s, err := convertUnicodeToString(s)

	if err != nil {
		return s, err
	}

	return s, nil
}

func (p *PL) PostalCode(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")

	s = p.rePostalCode.ReplaceAllString(s, "")

	if p.rePostalCodeIsNumeric.MatchString(s) {
		s = fmt.Sprintf("%s-%s", s[0:2], s[2:5])
	}

	return s, nil
}

func (p *PL) Street(s string) ([]string, error) {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "\n", "|")
	s = strings.ReplaceAll(s, ",", "|")
	s = replaceSpecialChars(s)

	s, err := convertUnicodeToString(s)

	if err != nil {
		return nil, err
	}

	s = p.reStreetDeleteWords.ReplaceAllString(s, "")
	s = p.reStreet.ReplaceAllString(s, "")
	s = p.reWordEndWithDot.ReplaceAllString(s, "")

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

func convertUnicodeToString(s string) (string, error) {

	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\"", "")
	result, err := strconv.Unquote("\"" + s + "\"")

	if err != nil {
		return result, err
	}

	return result, nil
}
