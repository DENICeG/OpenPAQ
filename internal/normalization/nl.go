// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type NL struct {
	rePostalCodeInvalidChars *regexp.Regexp
	rePostalCodeValid        *regexp.Regexp
	reCrap                   *regexp.Regexp
	reCity                   *regexp.Regexp
}

func newNL() (*NL, error) {
	rePostalCodeInvalidChars, errPostalCode := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;,-]")
	if errPostalCode != nil {
		return nil, errPostalCode
	}

	rePostalCodeValid, errPostalCode := regexp.Compile("\\b(\\d){4}\\s*([A-Z]){2}\\b")
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

	return &NL{
		rePostalCodeInvalidChars: rePostalCodeInvalidChars,
		rePostalCodeValid:        rePostalCodeValid,
		reCrap:                   reCrap,

		reCity: reCity,
	}, nil
}

func (nl *NL) GetCountryCode() string {
	return "nl"
}

func (nl *NL) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = nl.reCity.ReplaceAllString(s, "")

	return s, nil
}

func (nl *NL) PostalCode(s string) (string, error) {
	s = strings.ToUpper(s)
	s = nl.rePostalCodeInvalidChars.ReplaceAllString(s, "")

	if !nl.rePostalCodeValid.MatchString(s) {
		return s, errors.New("invalid postal code")
	}

	if len(s) == 6 {
		return fmt.Sprintf("%s %s", s[:4], s[4:]), nil
	}
	return s, nil
}

func (nl *NL) Street(s string) ([]string, error) {
	s = strings.ToLower(s)
	s = nl.reCrap.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\n", "|")
	s = strings.ReplaceAll(s, ",", "|")

	s = strings.ReplaceAll(s, "str.", "straat")
	s = strings.ReplaceAll(s, "ln.", "laan")
	s = strings.ReplaceAll(s, "wg.", "weg")
	s = strings.ReplaceAll(s, "pl.", "plein")
	s = strings.ReplaceAll(s, "gr.", "gracht")
	s = strings.ReplaceAll(s, "sgl.", "singel")
	s = strings.ReplaceAll(s, "kd.", "kade")
	s = strings.ReplaceAll(s, "hf.", "hof")

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

func (nl *NL) Nominatim(s string) (string, error) {
	return s, nil
}
