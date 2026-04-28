package normalization

import (
	"fmt"
	"regexp"
	"strings"
)

type AT struct {
	rePostalCode        *regexp.Regexp
	reCrap              *regexp.Regexp
	reStreetparts       *regexp.Regexp
	reStreetShortName   *regexp.Regexp
	reStreetHouseNumber *regexp.Regexp
	reCity              *regexp.Regexp
	reGNR               *regexp.Regexp
	reCityStopper       *regexp.Regexp
}

func newAT() (*AT, error) {
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
	reGNR, errGNR := regexp.Compile("\\b(gnr|eg|top|tuer|stock|stg|og|objekt|strases|gasse|stiege|road)\\b")
	if errGNR != nil {
		return nil, errGNR
	}
	reCityStopper, errCityStopper := regexp.Compile("\\b(sankt|im|am|ober|unter)\\b")
	if errCityStopper != nil {
		return nil, err
	}

	return &AT{
		rePostalCode:  rePostalCode,
		reCrap:        reCrap,
		reCity:        reCity,
		reGNR:         reGNR,
		reCityStopper: reCityStopper,
	}, nil
}

func (at *AT) GetCountryCode() string {
	return "at"
}

func (at *AT) City(s string) (string, error) {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = at.reCity.ReplaceAllString(s, "")
	s = at.reCityStopper.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "ü", "ue")
	s = strings.ReplaceAll(s, "ä", "ae")
	s = strings.ReplaceAll(s, "ö", "oe")
	s = strings.ReplaceAll(s, "ß", "ss")

	return s, nil
}

func (at *AT) PostalCode(s string) (string, error) {
	s = at.rePostalCode.ReplaceAllString(s, "")

	if len(s) != 4 {
		if len(s) > 4 {
			return s, fmt.Errorf("not valid postalcode")
		}
		return s, fmt.Errorf("not valid postalcode")
	}
	// remove leading 0 if present with 4 chars
	if len(s) == 4 && s[0] == '0' {
		return s[1:], nil
	}
	return s, nil
}

func (at *AT) Street(s string) ([]string, error) {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "\n", "|")
	s = strings.ReplaceAll(s, ",", "|")

	s = strings.ReplaceAll(s, "str.", "strasse")
	s = strings.ReplaceAll(s, "strasze", "strasse")
	s = strings.ReplaceAll(s, "chau.", "chaussee")
	s = strings.ReplaceAll(s, "rd.", "road")
	s = strings.ReplaceAll(s, "wr.", "wiener")
	s = at.reCrap.ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, "ü", "ue")
	s = strings.ReplaceAll(s, "ä", "ae")
	s = strings.ReplaceAll(s, "ö", "oe")
	s = strings.ReplaceAll(s, "ß", "ss")

	s = at.reGNR.ReplaceAllString(s, "")

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
