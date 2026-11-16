// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type DE struct {
	reCity           *regexp.Regexp
	rePostalCode     *regexp.Regexp
	reStreet         *regexp.Regexp
	reFoundStreet    *regexp.Regexp
	reStrAtEndOfWord *regexp.Regexp
}

func NewDE() (*DE, error) {
	reCity, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;0-9]")
	if err != nil {
		return nil, err
	}

	rePostalCode, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;a-zA-Z- ]")
	if err != nil {
		return nil, err
	}

	reStreet, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9]")
	if err != nil {
		return nil, err
	}

	reFoundStreet, err := regexp.Compile("([a-z]+straße)\\b")
	if err != nil {
		return nil, err
	}

	reStrAtEndOfWord, err := regexp.Compile("str\\b")
	if err != nil {
		return nil, err
	}

	return &DE{
		reCity:           reCity,
		rePostalCode:     rePostalCode,
		reStreet:         reStreet,
		reFoundStreet:    reFoundStreet,
		reStrAtEndOfWord: reStrAtEndOfWord,
	}, nil
}

func (g *DE) PostalCode(plz string) (string, error) {
	plz = g.rePostalCode.ReplaceAllString(plz, "")
	if len(plz) == 4 {
		plz = "0" + plz
	} else if len(plz) != 5 {
		return "", fmt.Errorf("PLZ does not containt 5 Characters")
	}
	if plz == "00000" {
		return "", fmt.Errorf("PLZ is 00000")
	}
	if _, err := strconv.Atoi(plz); err != nil {
		return "", err
	} else {
		return plz, nil
	}
}

func detectQuadratAddress(address string) (string, bool) {
	address = strings.ToLower(address)
	address = strings.ReplaceAll(address, ",", "\n")
	addressList := strings.Split(address, "\n")
	for _, item := range addressList {
		item = strings.TrimLeft(item, " ")
		if len(item) < 8 {
			if len(item) > 2 && unicode.IsLetter(rune(item[0])) && unicode.IsDigit(rune(item[1])) && string(item[2]) == " " {
				return item[0:2], true
			} else if len(item) == 3 && unicode.IsLetter(rune(item[0])) && string(item[1]) == " " && unicode.IsDigit(rune(item[2])) {
				return strings.ReplaceAll(item, " ", ""), true
			} else if len(item) == 2 && unicode.IsLetter(rune(item[0])) && unicode.IsDigit(rune(item[1])) {
				return item, true
			} else if len(item) > 3 && unicode.IsLetter(rune(item[0])) && string(item[1]) == " " && unicode.IsDigit(rune(item[2])) && string(item[3]) == " " {
				return strings.ReplaceAll(item, " ", "")[0:2], true
			}
		}
	}
	return "", false
}

func (g *DE) Street(address string) ([]string, error) {
	quadrat, isItQuadrat := detectQuadratAddress(address)

	address = strings.ToLower(address)

	address = strings.ReplaceAll(address, "/", " ")
	address = strings.ReplaceAll(address, "-", " ")
	address = strings.ReplaceAll(address, "\n", "|")
	address = strings.ReplaceAll(address, ",", "|")
	address = strings.ReplaceAll(address, "str.", "straße")
	address = strings.ReplaceAll(address, "strasse", "straße")
	address = strings.ReplaceAll(address, "strasze", "straße")
	address = strings.ReplaceAll(address, "chau.", "chaussee")
	address = strings.ReplaceAll(address, "rd.", "road")
	address = g.reStreet.ReplaceAllString(address, " ")

	address = strings.ReplaceAll(address, "ü", "ue")
	address = strings.ReplaceAll(address, "ä", "ae")
	address = strings.ReplaceAll(address, "ö", "oe")
	address = strings.ReplaceAll(address, "st.", "sankt")

	address = g.reStrAtEndOfWord.ReplaceAllString(address, "straße")

	foundStreets := g.reFoundStreet.Find([]byte(address))
	if len(foundStreets) != 0 {
		address = fmt.Sprintf("%s|%s", address, foundStreets)
	}

	addressSlice := strings.Split(address, "|")
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
	if isItQuadrat {
		cleanAddressSlice = append([]string{quadrat}, cleanAddressSlice...)
	}

	reducedCleanedAddresses := removeDuplicate(cleanAddressSlice)

	return reducedCleanedAddresses, nil
}

func (g *DE) City(city string) (string, error) {
	city = strings.ToLower(city)
	city = strings.ReplaceAll(city, "/", " ")
	city = strings.ReplaceAll(city, "-", " ")
	city = strings.ReplaceAll(city, "\n", " ")
	city = g.reCity.ReplaceAllString(city, " ")
	city = strings.ReplaceAll(city, "ü", "ue")
	city = strings.ReplaceAll(city, "ä", "ae")
	city = strings.ReplaceAll(city, "ö", "oe")
	city = strings.ReplaceAll(city, "st.", "sankt")
	city = strings.ReplaceAll(city, "ß", "sz")
	citySlice := strings.Split(city, " ")
	var cleanCitySlice []string
	for _, v := range citySlice {
		if len(v) > 1 {
			cleanCitySlice = append(cleanCitySlice, v)
		}
	}

	return strings.Join(cleanCitySlice, " "), nil
}

func (g *DE) GetCountryCode() string {
	return "de"
}
