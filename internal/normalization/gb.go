// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"errors"
	"regexp"
	"strings"
)

type GB struct {
	rePostalCodeInvalidChars      *regexp.Regexp
	rePostalCodeValidWithSpace    *regexp.Regexp
	rePostalCodeValidWithoutSpace *regexp.Regexp
	reStreetInvalidChars          *regexp.Regexp
	reStreetInvalidWords          *regexp.Regexp
	reStreetSplit                 *regexp.Regexp
	reStreetRd                    *regexp.Regexp
	reStreetDr                    *regexp.Regexp
	reStreetAve                   *regexp.Regexp
	reStreetSt                    *regexp.Regexp
	reStreetBlvd                  *regexp.Regexp
	reStreetHwy                   *regexp.Regexp
	reStreetPkwy                  *regexp.Regexp
	reCity                        *regexp.Regexp
}

func NewGB() (*GB, error) {

	rePostalCodeInvalidChars, err := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;\n\r-]")
	if err != nil {
		return nil, err
	}

	rePostalCodeValidWithSpace, err := regexp.Compile("\\b(\\w{1,2}((\\d\\w)|\\d{1,2})) (\\d\\w{2})\\b")
	if err != nil {
		return nil, err
	}

	rePostalCodeValidWithoutSpace, err := regexp.Compile("\\b(\\w{1,2}((\\d\\w)|\\d{1,2}))(\\d\\w{2})\\b")
	if err != nil {
		return nil, err
	}

	reCityInvalidChars, err := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;0-9\n\r]")
	if err != nil {
		return nil, err
	}

	reStreetInvalidChars, reStreetInvalidCharsErr := regexp.Compile("[+./(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9-\r]")
	if reStreetInvalidCharsErr != nil {
		return nil, reStreetInvalidCharsErr
	}

	reStreetInvalidWords, StreetInvalidWordsErr := regexp.Compile("\\b(floor|flat|unit|suite|fl|ste)\\b")
	if StreetInvalidWordsErr != nil {
		return nil, StreetInvalidWordsErr
	}

	reStreetSplit, StreetSplitErr := regexp.Compile("(\\d*[-|/]?\\d*) (\\s*|\\w*){1,3} (street|drive|road|avenue|boulevard|highway|parkway)")
	if StreetSplitErr != nil {
		return nil, StreetSplitErr
	}

	reStreetRd, StreetRdErr := regexp.Compile("\\brd\\b")
	if StreetRdErr != nil {
		return nil, StreetRdErr
	}

	reStreetDr, StreetDrErr := regexp.Compile("\\bdr\\b")
	if StreetDrErr != nil {
		return nil, StreetDrErr
	}

	reStreetAve, StreetAveErr := regexp.Compile("\\bave\\b")
	if StreetAveErr != nil {
		return nil, StreetAveErr
	}

	reStreetSt, StreetStErr := regexp.Compile("\\bst$")
	if StreetStErr != nil {
		return nil, StreetStErr
	}

	reStreetBlvd, StreetBlvdErr := regexp.Compile("\\bblvd\\b")
	if StreetBlvdErr != nil {
		return nil, StreetBlvdErr
	}

	reStreetHwy, StreetHwyErr := regexp.Compile("\\bhwy\\b")
	if StreetHwyErr != nil {
		return nil, StreetHwyErr
	}

	reStreetPkwy, StreetPkwyErr := regexp.Compile("\\bpkwy\\b")
	if StreetPkwyErr != nil {
		return nil, StreetPkwyErr
	}

	return &GB{
		rePostalCodeInvalidChars:      rePostalCodeInvalidChars,
		rePostalCodeValidWithSpace:    rePostalCodeValidWithSpace,
		rePostalCodeValidWithoutSpace: rePostalCodeValidWithoutSpace,
		reStreetInvalidChars:          reStreetInvalidChars,
		reStreetInvalidWords:          reStreetInvalidWords,
		reStreetSplit:                 reStreetSplit,
		reStreetRd:                    reStreetRd,
		reStreetDr:                    reStreetDr,
		reStreetAve:                   reStreetAve,
		reStreetSt:                    reStreetSt,
		reStreetBlvd:                  reStreetBlvd,
		reStreetHwy:                   reStreetHwy,
		reStreetPkwy:                  reStreetPkwy,
		reCity:                        reCityInvalidChars,
	}, nil
}

func (gb *GB) GetCountryCode() string {
	return "gb"
}

func (gb *GB) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = gb.reCity.ReplaceAllString(s, "")
	return strings.TrimSpace(s), nil
}

func (gb *GB) PostalCode(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = gb.rePostalCodeInvalidChars.ReplaceAllString(s, "")

	if !gb.rePostalCodeValidWithoutSpace.MatchString(s) {
		return "", errors.New("invalid postal code")
	}
	startlastThreeChars := len(s) - 3
	s = s[:startlastThreeChars]

	return s, nil
}

func (gb *GB) Street(s string) ([]string, error) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ",", "|")
	s = strings.ReplaceAll(s, "\n", "|")
	s = gb.rePostalCodeValidWithoutSpace.ReplaceAllString(s, "")
	s = gb.rePostalCodeValidWithSpace.ReplaceAllString(s, "")
	s = gb.reStreetSt.ReplaceAllString(s, "street")
	s = gb.reStreetDr.ReplaceAllString(s, "drive")
	s = gb.reStreetRd.ReplaceAllString(s, "road")
	s = gb.reStreetAve.ReplaceAllString(s, "avenue")
	s = gb.reStreetBlvd.ReplaceAllString(s, "boulevard")
	s = gb.reStreetHwy.ReplaceAllString(s, "highway")
	s = gb.reStreetPkwy.ReplaceAllString(s, "parkway")

	var tmpAddressSlice []string

	addressSlice := strings.Split(s, "|")

	for _, address := range addressSlice {
		found := gb.reStreetSplit.Find([]byte(address))
		if found != nil {
			foundStr := string(found)

			foundStr = gb.reStreetInvalidChars.ReplaceAllString(foundStr, "")
			foundStr = strings.TrimSpace(foundStr)
			tmpAddressSlice = append(tmpAddressSlice, foundStr)
		}
		address = gb.reStreetInvalidChars.ReplaceAllString(address, "")
		address = strings.TrimSpace(address)
		tmpAddressSlice = append(tmpAddressSlice, address)
	}

	addressSlice = tmpAddressSlice

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
