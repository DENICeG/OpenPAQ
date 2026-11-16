// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"regexp"
	"strings"
)

type US struct {
	rePostalCode          *regexp.Regexp
	reSpecialCharacters   *regexp.Regexp
	reNumbers             *regexp.Regexp
	reNonNumbers          *regexp.Regexp
	reNorth               *regexp.Regexp
	reSouth               *regexp.Regexp
	reWest                *regexp.Regexp
	reEast                *regexp.Regexp
	reNorthEast           *regexp.Regexp
	reNorthWest           *regexp.Regexp
	reSouthEast           *regexp.Regexp
	reSouthWest           *regexp.Regexp
	reStreetFilter        *regexp.Regexp
	reStreetRd            *regexp.Regexp
	reStreetDr            *regexp.Regexp
	reStreetAve           *regexp.Regexp
	reStreetSt            *regexp.Regexp
	reStreetBlvd          *regexp.Regexp
	reStreetHwy           *regexp.Regexp
	reStreetPkwy          *regexp.Regexp
	reStreetNumberPersist *regexp.Regexp
}

func newUS() (*US, error) {

	rePostalCode, errPostalCode := regexp.Compile("-.*|[a-zA-Z]|[^0-9]")
	if errPostalCode != nil {
		return nil, errPostalCode
	}

	reSpecialCharacters, errSpecialCharacters := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\".:;]")
	if errSpecialCharacters != nil {
		return nil, errSpecialCharacters
	}

	reNumbers, errNumbers := regexp.Compile("(\\d+)")
	if errNumbers != nil {
		return nil, errNumbers
	}

	reNonNumbers, errNonNumbers := regexp.Compile("(\\D+)")
	if errNonNumbers != nil {
		return nil, errNonNumbers
	}

	reNorth, errNorth := regexp.Compile("\\bn\\b")
	if errNorth != nil {
		return nil, errNorth
	}

	reSouth, errSouth := regexp.Compile("\\bs\\b")
	if errSouth != nil {
		return nil, errSouth
	}

	reWest, errWest := regexp.Compile("\\bw\\b")
	if errWest != nil {
		return nil, errWest
	}

	reEast, errEast := regexp.Compile("\\be\\b")
	if errEast != nil {
		return nil, errEast
	}

	reNorthEast, errNE := regexp.Compile("\\bne\\b")
	if errNE != nil {
		return nil, errNE
	}

	reNorthWest, errNW := regexp.Compile("\\bnw\\b")
	if errNW != nil {
		return nil, errNW
	}

	reSouthEast, errSE := regexp.Compile("\\bse\\b")
	if errSE != nil {
		return nil, errSE
	}

	reSouthWest, errSW := regexp.Compile("\\bsw\\b")
	if errSW != nil {
		return nil, errSW
	}

	reStreetFilter, StreetFilter := regexp.Compile("\\b(floor|suite|fl|ste)\\b")
	if StreetFilter != nil {
		return nil, StreetFilter
	}

	reStreetRd, StreetStreetRd := regexp.Compile("\\brd\\b")
	if StreetStreetRd != nil {
		return nil, StreetStreetRd
	}

	reStreetDr, StreetStreetDr := regexp.Compile("\\bdr\\b")
	if StreetStreetDr != nil {
		return nil, StreetStreetDr
	}

	reStreetAve, StreetStreetAve := regexp.Compile("\\bave\\b")
	if StreetStreetAve != nil {
		return nil, StreetStreetAve
	}

	reStreetSt, StreetStreetSt := regexp.Compile("\\bst\\b")
	if StreetStreetSt != nil {
		return nil, StreetStreetSt
	}

	reStreetBlvd, StreetStreetBlvd := regexp.Compile("\\bblvd\\b")
	if StreetStreetBlvd != nil {
		return nil, StreetStreetBlvd
	}

	reStreetHwy, StreetStreetHwy := regexp.Compile("\\bhwy\\b")
	if StreetStreetHwy != nil {
		return nil, StreetStreetHwy
	}

	reStreetPkwy, StreetStreetPkwy := regexp.Compile("\\bpkwy\\b")
	if StreetStreetPkwy != nil {
		return nil, StreetStreetPkwy
	}

	reStreetNumberPersist, StreetNumberPersist := regexp.Compile("(\\d+\\W?)(st|nd|rd|th)\\W?(avenue|road|street)")
	if StreetNumberPersist != nil {
		return nil, StreetNumberPersist
	}

	return &US{
		rePostalCode:          rePostalCode,
		reSpecialCharacters:   reSpecialCharacters,
		reNumbers:             reNumbers,
		reNonNumbers:          reNonNumbers,
		reNorth:               reNorth,
		reSouth:               reSouth,
		reWest:                reWest,
		reEast:                reEast,
		reNorthEast:           reNorthEast,
		reNorthWest:           reNorthWest,
		reSouthEast:           reSouthEast,
		reSouthWest:           reSouthWest,
		reStreetFilter:        reStreetFilter,
		reStreetRd:            reStreetRd,
		reStreetDr:            reStreetDr,
		reStreetAve:           reStreetAve,
		reStreetSt:            reStreetSt,
		reStreetBlvd:          reStreetBlvd,
		reStreetHwy:           reStreetHwy,
		reStreetPkwy:          reStreetPkwy,
		reStreetNumberPersist: reStreetNumberPersist,
	}, nil
}

func (us *US) GetCountryCode() string {
	return "us"
}

func (us *US) City(s string) (string, error) {

	normalizedCity := strings.ToLower(s)
	normalizedCity = strings.TrimSpace(normalizedCity)
	normalizedCity = us.reSpecialCharacters.ReplaceAllString(normalizedCity, "")
	normalizedCity = us.reNumbers.ReplaceAllString(normalizedCity, "")
	return normalizedCity, nil
}

func (us *US) PostalCode(s string) (string, error) {

	normalizerPostalCode := us.rePostalCode.ReplaceAllString(s, "")
	normalizerPostalCode = us.reNonNumbers.ReplaceAllString(normalizerPostalCode, "")

	return normalizerPostalCode, nil
}

func (us *US) Street(s string) ([]string, error) {

	normalizedStreet := strings.ToLower(s)
	normalizedStreet = strings.TrimSpace(normalizedStreet)

	normalizedStreet = strings.ReplaceAll(normalizedStreet, "\n", "|")
	normalizedStreet = strings.ReplaceAll(normalizedStreet, ",", "|")

	normalizedStreet = us.reSpecialCharacters.ReplaceAllString(normalizedStreet, "")

	normalizedStreet = us.reNorth.ReplaceAllString(normalizedStreet, "north")
	normalizedStreet = us.reSouth.ReplaceAllString(normalizedStreet, "south")
	normalizedStreet = us.reWest.ReplaceAllString(normalizedStreet, "west")
	normalizedStreet = us.reEast.ReplaceAllString(normalizedStreet, "east")
	normalizedStreet = us.reNorthEast.ReplaceAllString(normalizedStreet, "northeast")
	normalizedStreet = us.reNorthWest.ReplaceAllString(normalizedStreet, "northwest")
	normalizedStreet = us.reSouthEast.ReplaceAllString(normalizedStreet, "southeast")
	normalizedStreet = us.reSouthWest.ReplaceAllString(normalizedStreet, "southwest")

	normalizedStreet = us.reStreetFilter.ReplaceAllString(normalizedStreet, "")

	normalizedStreet = us.reStreetSt.ReplaceAllString(normalizedStreet, "street")
	normalizedStreet = us.reStreetDr.ReplaceAllString(normalizedStreet, "drive")
	normalizedStreet = us.reStreetRd.ReplaceAllString(normalizedStreet, "road")
	normalizedStreet = us.reStreetAve.ReplaceAllString(normalizedStreet, "avenue")
	normalizedStreet = us.reStreetBlvd.ReplaceAllString(normalizedStreet, "boulevard")
	normalizedStreet = us.reStreetHwy.ReplaceAllString(normalizedStreet, "highway")
	normalizedStreet = us.reStreetPkwy.ReplaceAllString(normalizedStreet, "parkway")

	numberPersistIndex := us.reStreetNumberPersist.FindAllIndex(([]byte)(normalizedStreet), -1)
	numberIndex := us.reNumbers.FindAllIndex(([]byte)(normalizedStreet), -1)

	tempNormalizedStreet := ""
	tempKeepStart := 0

	if len(numberIndex) > 0 {
		for numDropIdx, numDropValue := range numberIndex {
			tempKeepEnd := numDropValue[0]

			if tempKeepEnd > tempKeepStart {
				tempNormalizedStreet = tempNormalizedStreet + normalizedStreet[tempKeepStart:tempKeepEnd]
			}

			if len(numberPersistIndex) == 0 {
				tempKeepStart = numDropValue[1]
			}

			for _, numIndexPersist := range numberPersistIndex {
				if numIndexPersist[0] == numDropValue[0] {
					tempKeepStart = numIndexPersist[0]
					break
				}
				tempKeepStart = numDropValue[1]
			}

			if numDropIdx == len(numberIndex)-1 {
				tempNormalizedStreet = tempNormalizedStreet + normalizedStreet[tempKeepStart:]
			}
		}
		normalizedStreet = tempNormalizedStreet
	}

	addressSlice := strings.Split(normalizedStreet, "|")
	var cleanAddressSlice []string
	for _, v := range addressSlice {
		var cleanedParts []string
		addressPartSlices := strings.Split(v, " ")
		for _, p := range addressPartSlices {
			if len(p) > 2 {
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
