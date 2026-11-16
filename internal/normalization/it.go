// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type IT struct {
	reStreetInvalidChars     *regexp.Regexp
	reStreetContainsVia      *regexp.Regexp
	rePostalCodeInvalidChars *regexp.Regexp
	rePostalCodeValid        *regexp.Regexp
	reCorso                  *regexp.Regexp
	rePiazza                 *regexp.Regexp
	reViale                  *regexp.Regexp
	reVia                    *regexp.Regexp
	reVicolo                 *regexp.Regexp
	reLargo                  *regexp.Regexp
	reTraversa               *regexp.Regexp
	reCityInvalidChars       *regexp.Regexp
}

func NewIT() (*IT, error) {
	reStreetInvalidChars, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9\n\r]")
	if err != nil {
		return nil, err
	}

	reStreetContainsVia, err := regexp.Compile("\\b(via|piazza|viale)( [a-z]*[.]?[a-z]*)*")
	if err != nil {
		return nil, err
	}

	rePostalCodeInvalidChars, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;\n\ra-z]")
	if err != nil {
		return nil, err
	}

	rePostalCodeValid, err := regexp.Compile("\\b\\d{5}\\b")
	if err != nil {
		return nil, err
	}

	reCityInvalidChars, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;0-9\n\r]")
	if err != nil {
		return nil, err
	}

	reCorso, err := regexp.Compile("\\bc\\.so\\b")
	if err != nil {
		return nil, err
	}

	rePiazza, err := regexp.Compile("\\bp\\.za\\b")
	if err != nil {
		return nil, err
	}

	reViale, err := regexp.Compile("\\bv\\.le\\b")
	if err != nil {
		return nil, err
	}

	reVia, err := regexp.Compile("\\bv\\.")
	if err != nil {
		return nil, err
	}

	reVicolo, err := regexp.Compile("\\b\"v\\.lo\\b")
	if err != nil {
		return nil, err
	}

	reLargo, err := regexp.Compile("\\bl\\.go\\b")
	if err != nil {
		return nil, err
	}

	reTraversa, err := regexp.Compile("\\btrav\\.")
	if err != nil {
		return nil, err
	}

	return &IT{
		reStreetInvalidChars:     reStreetInvalidChars,
		reStreetContainsVia:      reStreetContainsVia,
		rePostalCodeInvalidChars: rePostalCodeInvalidChars,
		rePostalCodeValid:        rePostalCodeValid,
		reCorso:                  reCorso,
		rePiazza:                 rePiazza,
		reViale:                  reViale,
		reVia:                    reVia,
		reVicolo:                 reVicolo,
		reLargo:                  reLargo,
		reTraversa:               reTraversa,
		reCityInvalidChars:       reCityInvalidChars,
	}, nil
}

func (g *IT) GetCountryCode() string {
	return "it"
}

func (g *IT) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = g.reCityInvalidChars.ReplaceAllString(s, "")
	return strings.TrimSpace(s), nil
}

func (g *IT) PostalCode(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = g.rePostalCodeInvalidChars.ReplaceAllString(s, "")

	if !g.rePostalCodeValid.MatchString(s) {
		return s, errors.New("invalid postal code")
	}

	return s, nil

}

func (g *IT) Street(s string) ([]string, error) {

	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ",", "|")
	s = strings.ReplaceAll(s, "\n", "|")

	s = g.reCorso.ReplaceAllString(s, "corso")
	s = g.rePiazza.ReplaceAllString(s, "piazza")
	s = g.reViale.ReplaceAllString(s, "viale")
	s = g.reVia.ReplaceAllString(s, "via")
	s = g.reVicolo.ReplaceAllString(s, "vicolo")
	s = g.reLargo.ReplaceAllString(s, "largo")
	s = g.reTraversa.ReplaceAllString(s, "traversa")

	s = g.reStreetInvalidChars.ReplaceAllString(s, "")

	foundStreets := g.reStreetContainsVia.Find([]byte(s))
	if len(foundStreets) != 0 {
		s = fmt.Sprintf("%s|%s", s, foundStreets)

	}

	s = replaceItalianShortcuts(s)

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

	reducedCleanedAddresses := removeDuplicate(cleanAddressSlice)

	return reducedCleanedAddresses, nil

}

func replaceItalianShortcuts(s string) string {

	result := s

	replacments := map[string][]string{
		" g.":  {" giovanni", " giuseppe", " giacomo", " gabriele", " giorgio"},
		"str.": {"strada", "strasse"},
		" st.": {" strada", " santo", " santa"},
		" f.":  {" francesco", " filippo", " ferdinando"},
		" s.":  {" santo", " santa"},
		" v.":  {" vittorio"},
		" m.":  {" marco", " maria", " michele"},
		" d.":  {" don"},
	}

	for key, value := range replacments {
		if strings.Contains(s, key) {
			for _, v := range value {
				temp := strings.ReplaceAll(s, key, v)
				result = fmt.Sprintf("%s|%s", result, temp)

			}
		}
	}

	return result

}
