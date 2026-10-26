package normalization

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type DK struct {
	reCompanySuffix       *regexp.Regexp // can contain c/o information and should be ignored
	reBuildingAfterNumber *regexp.Regexp // can contain additional location info, such as "entrance", "building" etc.
	reFloor               *regexp.Regexp
	reDoor                *regexp.Regexp
	reCrap                *regexp.Regexp
	rePostalCode          *regexp.Regexp
}

// Only these known postal districts may be reduced to their base city for
// matching. Other suffixes can distinguish different towns (Nykøbing F/M/S,
// for example) and must be preserved. This is not a postcode validation table.
// District names: https://www.statistikbanken.dk/POSTNR1
var dkCityDistrictSuffixes = map[string][]string{
	"københavn":     {"k", "v", "n", "ø", "s", "nv", "sv"},
	"frederiksberg": {"c"},
	"aarhus":        {"c", "n", "v"},
	"aalborg":       {"sv", "sø", "øst"},
	"odense":        {"c", "m", "n", "nv", "nø", "s", "sv", "sø", "v"},
	"esbjerg":       {"n", "v", "ø"},
	"randers":       {"c", "nv", "nø", "sv", "sø"},
}

func newDK() (*DK, error) {
	reCompanySuffix, errCompanySuffix := regexp.Compile("\\s*( v/| c/o| att)")
	if errCompanySuffix != nil {
		return nil, errCompanySuffix
	}
	reBuilding, errBuildings := regexp.Compile("\\W(sal|etage|floor|kl|kld|kælder|st|stuen|parterre|dør|port|opg|opgang|indgang|bygning|tv|th|bygn|mf)\\b")
	if errBuildings != nil {
		return nil, errBuildings
	}

	reFloor, errFloor := regexp.Compile("(\\W)\\b(kld|kl|kælder\\w*|st|stuen|parterren)\\b")
	if errFloor != nil {
		return nil, errFloor
	}

	reDoor, errDoor := regexp.Compile("(\\d.*)(?:\\btv|\\bt\\W+v|v\\.|\\bv|\\bth|\\bt\\W+h|h\\.|\\bh|mf)\\b")
	if errDoor != nil {
		return nil, errDoor
	}

	rePostalCode, errPostalCode := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;a-zA-Z- ]")
	if errPostalCode != nil {
		return nil, errPostalCode
	}

	reCrap, errCrap := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9]")
	if errCrap != nil {
		return nil, errCrap
	}

	return &DK{
		reCompanySuffix:       reCompanySuffix,
		reBuildingAfterNumber: reBuilding,
		reFloor:               reFloor,
		reDoor:                reDoor,
		rePostalCode:          rePostalCode,
		reCrap:                reCrap,
	}, nil
}

func (dk *DK) GetCountryCode() string {
	return "dk"
}

func (dk *DK) City(s string) (string, error) {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "oe", "ø")
	s = strings.ReplaceAll(s, "ae", "æ")
	s = strings.Join(strings.Fields(s), " ")

	// City returns the comparison name; the original city remains in types.Input.
	// Match a whole trailing token so names such as Hvide Sande are not damaged.
	if i := strings.LastIndexByte(s, ' '); i >= 0 {
		base := normalizeDKCityName(s[:i])
		suffix := strings.TrimSuffix(s[i+1:], ".")
		if slices.Contains(dkCityDistrictSuffixes[base], suffix) {
			return base, nil
		}
	}

	return normalizeDKCityName(s), nil
}

func normalizeDKCityName(s string) string {
	switch s {
	case "kbh", "kbh.", "kobenhavn", "copenhagen", "kopenhagen":
		return "københavn"
	case "århus":
		return "aarhus"
	}
	return s
}

func (dk *DK) PostalCode(s string) (string, error) {
	s = dk.rePostalCode.ReplaceAllString(s, "")

	if len(s) != 4 {
		if len(s) > 4 {
			return s[:4], fmt.Errorf("not valid postalcode")
		}
		return s, fmt.Errorf("not valid postalcode")
	}

	return s, nil
}

func (dk *DK) Street(s string) ([]string, error) {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "\n", "|")
	s = strings.ReplaceAll(s, ",", "|")

	s = strings.ReplaceAll(s, "oe", "ø")
	s = strings.ReplaceAll(s, "ae", "æ")

	s = dk.reDoor.ReplaceAllString(s, " ")
	s = dk.reCompanySuffix.ReplaceAllString(s, " ")
	s = dk.reBuildingAfterNumber.ReplaceAllString(s, " ")
	s = dk.reFloor.ReplaceAllString(s, " ")

	s = dk.reCrap.ReplaceAllString(s, " ")

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
