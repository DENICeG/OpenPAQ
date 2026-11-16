// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"fmt"
	"regexp"
	"strings"
)

type fr struct {
	reStreet      *regexp.Regexp
	rePostalCode  *regexp.Regexp
	reCity        *regexp.Regexp
	reBoulevard   *regexp.Regexp
	reAvenue      *regexp.Regexp
	reBis         *regexp.Regexp
	reZac         *regexp.Regexp
	reFoundStreet *regexp.Regexp
}

func newFR() (*fr, error) {
	reStreet, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\":;0-9\n\r]")
	if err != nil {
		return nil, err
	}

	rePostalCode, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;\n\ra-zA-Z]")
	if err != nil {
		return nil, err
	}

	reCity, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;0-9\n\r]")
	if err != nil {
		return nil, err
	}

	reBoulevard, err := regexp.Compile("\\b(bd|bld)\\b[.]?")
	if err != nil {
		return nil, err
	}

	reAvenue, err := regexp.Compile("\\bave\\b")
	if err != nil {
		return nil, err
	}

	//https://www.reddit.com/r/French/comments/40d6hj/the_word_bis_in_frenchswiss_addresses/
	reBis, err := regexp.Compile("\\bbis\\b")
	if err != nil {
		return nil, err
	}

	reZac, err := regexp.Compile("\\bzac\\b")
	if err != nil {
		return nil, err
	}

	reFoundStreet, err := regexp.Compile("\\b(rue|allee|avenue|boulevard)( [a-z]*[.|']?[a-z]*)*")
	if err != nil {
		return nil, err
	}

	return &fr{
		reStreet:      reStreet,
		rePostalCode:  rePostalCode,
		reCity:        reCity,
		reBoulevard:   reBoulevard,
		reAvenue:      reAvenue,
		reBis:         reBis,
		reZac:         reZac,
		reFoundStreet: reFoundStreet,
	}, nil

}

func (g *fr) GetCountryCode() string {
	return "fr"
}

func (g *fr) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = g.reCity.ReplaceAllString(s, "")
	s = replaceFrLetters(s)
	s = strings.TrimSpace(s)
	return s, nil
}

func (g *fr) PostalCode(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = g.rePostalCode.ReplaceAllString(s, "")
	return s, nil
}

func (g *fr) Street(s string) ([]string, error) {

	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ",", "|")
	s = strings.ReplaceAll(s, "\n", "|")
	s = strings.ReplaceAll(s, " - ", "|")
	s = g.reStreet.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "- ", "")
	s = strings.ReplaceAll(s, " -", "")
	s = replaceFrLetters(s)
	s = g.replaceShortcuts(s)
	s = g.reZac.ReplaceAllString(s, "|zac")
	s = strings.TrimSpace(s)

	foundStreets := g.reFoundStreet.Find([]byte(s))
	if len(foundStreets) != 0 {
		s = fmt.Sprintf("%s|%s", s, foundStreets)

	}

	addressSlice := strings.Split(s, "|")
	var cleanAddressSlice []string
	for _, v := range addressSlice {
		v = strings.TrimSpace(v)
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

	reducedCleanedAddresses := removeDuplicate(cleanAddressSlice)

	return reducedCleanedAddresses, nil

}
func (g *fr) replaceShortcuts(s string) string {
	s = g.reBoulevard.ReplaceAllString(s, "boulevard")
	s = g.reAvenue.ReplaceAllString(s, "avenue")
	s = g.reBis.ReplaceAllString(s, "")
	return s

}

func replaceFrLetters(s string) string {
	s = strings.ReplaceAll(s, "í", "i")
	s = strings.ReplaceAll(s, "é", "e")
	s = strings.ReplaceAll(s, "è", "e")
	s = strings.ReplaceAll(s, "ê", "e")
	s = strings.ReplaceAll(s, "ë", "e")
	s = strings.ReplaceAll(s, "à", "a")
	s = strings.ReplaceAll(s, "â", "a")
	s = strings.ReplaceAll(s, "ä", "a")
	s = strings.ReplaceAll(s, "ù", "u")
	s = strings.ReplaceAll(s, "û", "u")
	s = strings.ReplaceAll(s, "ü", "u")
	s = strings.ReplaceAll(s, "ô", "o")
	s = strings.ReplaceAll(s, "ö", "o")
	s = strings.ReplaceAll(s, "ç", "c")
	s = strings.ReplaceAll(s, "œ", "oe")

	return s
}
