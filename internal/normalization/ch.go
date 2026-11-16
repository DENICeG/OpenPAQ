// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"regexp"
	"strings"
)

type CH struct {
	reStreet     *regexp.Regexp
	rePostalCode *regexp.Regexp
	reCity       *regexp.Regexp
}

func NewCh() (*CH, error) {
	reStreet, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9\n\r]")
	if err != nil {
		return nil, err
	}

	rePostalCode, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_a-z\":;\n\r]")
	if err != nil {
		return nil, err
	}

	reCity, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;0-9\n\r]")
	if err != nil {
		return nil, err
	}
	return &CH{
		reStreet:     reStreet,
		rePostalCode: rePostalCode,
		reCity:       reCity,
	}, nil
}

func (ch *CH) GetCountryCode() string {
	return "ch"
}

func (ch *CH) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = ch.reCity.ReplaceAllString(s, "")
	s = replaceChLetters(s)
	return strings.TrimSpace(s), nil
}

func (ch *CH) PostalCode(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	return ch.rePostalCode.ReplaceAllString(s, ""), nil
}

func (ch *CH) Street(s string) ([]string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ",", "|")
	s = strings.ReplaceAll(s, "\n", "|")
	s = strings.ReplaceAll(s, "str.", "straße")
	s = strings.ReplaceAll(s, "strasse", "straße")
	s = ch.reStreet.ReplaceAllString(s, "")
	s = replaceChLetters(s)

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

func replaceChLetters(s string) string {
	s = strings.ReplaceAll(s, "í", "i")
	s = strings.ReplaceAll(s, "é", "e")
	s = strings.ReplaceAll(s, "è", "e")
	s = strings.ReplaceAll(s, "ê", "e")
	s = strings.ReplaceAll(s, "ë", "e")
	s = strings.ReplaceAll(s, "à", "a")
	s = strings.ReplaceAll(s, "â", "a")
	s = strings.ReplaceAll(s, "ä", "ae")
	s = strings.ReplaceAll(s, "ù", "u")
	s = strings.ReplaceAll(s, "û", "u")
	s = strings.ReplaceAll(s, "ü", "ue")
	s = strings.ReplaceAll(s, "ô", "o")
	s = strings.ReplaceAll(s, "ö", "oe")
	s = strings.ReplaceAll(s, "ç", "c")
	s = strings.ReplaceAll(s, "œ", "oe")

	return s
}
