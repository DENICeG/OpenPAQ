package normalization

import (
	"slices"
	"testing"
)

func TestGetCountryCode(t *testing.T) {
	dkNormalizer, _ := newDK()

	countryCode := dkNormalizer.GetCountryCode()

	if countryCode != "dk" {
		t.Errorf("got %s, want dk", countryCode)
	}
}

func TestCity(t *testing.T) {
	dkNormalizer, err := newDK()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Want only small letters", input: "ACityWithCapitalLetters", expect: "acitywithcapitalletters"},
		{name: "oe should be replaced", input: "ACityWithOeWithin", expect: "acitywithøwithin"},
		{name: "ae should be replaced", input: "ACityWithAeWithin", expect: "acitywithæwithin"},
		{name: "Copenhagen centre", input: "København K", expect: "københavn"},
		{name: "Copenhagen west", input: "København V", expect: "københavn"},
		{name: "Copenhagen north", input: "København N", expect: "københavn"},
		{name: "Copenhagen east", input: "København Ø", expect: "københavn"},
		{name: "Copenhagen south", input: "København S", expect: "københavn"},
		{name: "Copenhagen northwest", input: "København NV", expect: "københavn"},
		{name: "Copenhagen southwest", input: "København SV", expect: "københavn"},
		{name: "Frederiksberg centre", input: "Frederiksberg C", expect: "frederiksberg"},
		{name: "Aarhus centre", input: "Aarhus C", expect: "aarhus"},
		{name: "Aarhus north", input: "Aarhus N", expect: "aarhus"},
		{name: "Aarhus west", input: "Aarhus V", expect: "aarhus"},
		{name: "Aalborg southeast", input: "Aalborg SØ", expect: "aalborg"},
		{name: "Aalborg southwest", input: "Aalborg SV", expect: "aalborg"},
		{name: "Aalborg east is written in full", input: "Aalborg Øst", expect: "aalborg"},
		{name: "Odense northeast", input: "Odense NØ", expect: "odense"},
		{name: "Odense M is a district of the same city", input: "Odense M", expect: "odense"},
		{name: "Esbjerg east", input: "Esbjerg Ø", expect: "esbjerg"},
		{name: "Randers northwest", input: "Randers NV", expect: "randers"},
		{name: "ASCII spelling with district", input: "Koebenhavn K", expect: "københavn"},
		{name: "simplified spelling with district", input: "Kobenhavn K", expect: "københavn"},
		{name: "English spelling with district", input: "Copenhagen K", expect: "københavn"},
		{name: "German spelling with district", input: "Kopenhagen K", expect: "københavn"},
		{name: "abbreviated name with district", input: "Kbh. K", expect: "københavn"},
		{name: "abbreviated name without district", input: "Kbh", expect: "københavn"},
		{name: "historic spelling with district", input: "Århus C", expect: "aarhus"},
		{name: "historic spelling without district", input: "Århus", expect: "aarhus"},
		{name: "whitespace and trailing suffix period", input: " \tKØBENHAVN\u00a0  K.\n", expect: "københavn"},
		{name: "unknown city keeps directional suffix", input: "City SV", expect: "city sv"},
		{name: "unknown district is preserved", input: "København X", expect: "københavn x"},
		{name: "district must belong to this city", input: "Aarhus K", expect: "aarhus k"},
		{name: "district token must be complete", input: "København Nord", expect: "københavn nord"},
		{name: "suffix must be last", input: "København K Extra", expect: "københavn k extra"},
		{name: "attached suffix is preserved", input: "KøbenhavnK", expect: "københavnk"},
		{name: "Falster must remain distinct", input: "Nykøbing F", expect: "nykøbing f"},
		{name: "Mors must remain distinct", input: "Nykøbing M", expect: "nykøbing m"},
		{name: "Sjaelland must remain distinct", input: "Nykøbing S", expect: "nykøbing s"},
		{name: "Viby J suffix is preserved", input: "Viby J", expect: "viby j"},
		{name: "word beginning with s is preserved", input: "Hvide Sande", expect: "hvide sande"},
		{name: "word beginning with v is preserved", input: "Sønder Vissing", expect: "sønder vissing"},
		{name: "word beginning with n is preserved", input: "Store Nørreby", expect: "store nørreby"},
		{name: "multiword city whitespace", input: "  Nørre\tSnede  ", expect: "nørre snede"},
		{name: "hyphenated name is preserved", input: "Nørre-Åby", expect: "nørre-åby"},
		{name: "empty city", input: "", expect: ""},
		{name: "whitespace only", input: " \t\n", expect: ""},
		{name: "sv should not be removed at the beginning", input: "SV City", expect: "sv city"},
		{name: "sv should not be removed at the beginning", input: "SVCity", expect: "svcity"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := dkNormalizer.City(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeNoError(t *testing.T) {
	dkNormalizer, _ := newDK()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be removed", input: "1234abcdefghijklmnopqrstuvwxyz", expect: "1234"},
		{name: "Capital letters should be removed", input: "1234ABCDEFGHIJKLMNOPQRSTUVWXYZ", expect: "1234"},
		{name: "Special characters should be removed", input: "1234+/(){}[]<>!§$%&=?*#€¿_\",:;-", expect: "1234"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := dkNormalizer.PostalCode(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeExpectError(t *testing.T) {
	dkNormalizer, _ := newDK()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Postal code to long", input: "12345", expect: "1234"},
		{name: "Postal code to short", input: "123", expect: "123"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := dkNormalizer.PostalCode(test.input)
			if err == nil {
				t.Error("expect error got none")
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestStreet(t *testing.T) {
	dkNormalizer, _ := newDK()

	tests := []struct {
		name            string
		input           string
		expectedContent []string
	}{
		{name: "no separator", input: "JustAStreet", expectedContent: []string{"justastreet"}},
		{name: "spaces", input: "Just A Street", expectedContent: []string{"just street"}},
		{name: "comma", input: "Just , Street", expectedContent: []string{"just", "street"}},
		{name: "new line", input: "Just \n Street", expectedContent: []string{"just", "street"}},
		{name: "oe", input: "Just Streetoe", expectedContent: []string{"just streetø"}},
		{name: "ae", input: "Just Streetae", expectedContent: []string{"just streetæ"}},
		{name: "remove v at the end", input: "Just Street v", expectedContent: []string{"just street"}},
		{name: "remove v at the end of word", input: "Just Streetv", expectedContent: []string{"just streetv"}},
		{name: "c/o at the end", input: "Just Street c/o", expectedContent: []string{"just street"}},
		{name: "att at the end", input: "Just Street att", expectedContent: []string{"just street"}},
		{name: "att at the end of word", input: "Just Streetatt", expectedContent: []string{"just streetatt"}},
		{name: "att at the end of word", input: "Just Streetatt", expectedContent: []string{"just streetatt"}},
		{name: "building context", input: "Just Street sal etage floor kl kld kælder st stuen parterre dør port opg opgang indgang bygning bygn tv th mf", expectedContent: []string{"just street"}},
		{name: "door context test-1", input: "Test Vænge 9, opgang 4, 2. tv blub", expectedContent: []string{"test vænge blub"}},
		{name: "door context test-2", input: "lurchschlumpf 2. tv asdf", expectedContent: []string{"lurchschlumpf asdf"}},
		{name: "door context test-3", input: "lurchschlumpf 1.tv.", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-4", input: "lurchschlumpf 1.t.v.", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-5", input: "lurchschlumpf 1.v.", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-6", input: "lurchschlumpf 1. v.", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-7", input: "lurchschlumpf 1. th", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-8", input: "lurchschlumpf 1. h", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-9", input: "lurchschlumpf 1 h", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-10", input: "lurchschlumpf 1 v", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-11", input: "lurchschlumpf 1. mf", expectedContent: []string{"lurchschlumpf"}},
		{name: "floor context", input: "lurchschlumpf kld kl kælder st stuen parterren", expectedContent: []string{"lurchschlumpf"}},
		{name: "floor context", input: "lurchschlumpfkld", expectedContent: []string{"lurchschlumpfkld"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := dkNormalizer.Street(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if !slices.Equal(test.expectedContent, res) {
				t.Errorf("got %v, want %v", res, test.expectedContent)
			}
		})
	}
}
