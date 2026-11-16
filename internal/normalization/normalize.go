// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

type Normalize interface {
	GetCountryCode() string
	City(string) (string, error)
	PostalCode(string) (string, error)
	Street(string) ([]string, error)
}

type Normalizer struct {
	availableCountries []Normalize
	fallback           Normalize
}

func NewNormalizer(fallback string) *Normalizer {
	n := &Normalizer{}

	genericNormalizer, err := NewGeneric()
	if err != nil {
		panic(err)
	}

	deNormalizer, err := NewDE()
	if err != nil {
		panic(err)
	}

	dkNormalizer, err := newDK()
	if err != nil {
		panic(err)
	}

	usNormalizer, err := newUS()
	if err != nil {
		panic(err)
	}

	plNormalizer, err := NewPl()
	if err != nil {
		panic(err)
	}

	atNormalizer, err := newAT()
	if err != nil {
		panic(err)
	}

	ukNormalizer, err := NewGB()
	if err != nil {
		panic(err)
	}

	esNormalizer, err := newES()
	if err != nil {
		panic(err)
	}

	itNormalizer, err := NewIT()
	if err != nil {
		panic(err)
	}

	nlNormalizer, err := newNL()
	if err != nil {
		panic(err)
	}

	frNormalizer, err := newFR()
	if err != nil {
		panic(err)
	}

	chNormalizer, err := NewCh()
	if err != nil {
		panic(err)
	}

	n.register(genericNormalizer)
	n.register(deNormalizer)
	n.register(dkNormalizer)
	n.register(usNormalizer)
	n.register(esNormalizer)
	n.register(plNormalizer)
	n.register(atNormalizer)
	n.register(ukNormalizer)
	n.register(nlNormalizer)
	n.register(itNormalizer)
	n.register(frNormalizer)
	n.register(chNormalizer)

	for i := range n.availableCountries {
		if fallback == n.availableCountries[i].GetCountryCode() {
			n.fallback = n.availableCountries[i]
		}
	}

	return n
}

func (n *Normalizer) register(n2 Normalize) {
	for _, v := range n.availableCountries {
		if v.GetCountryCode() == n2.GetCountryCode() {
			return
		}
	}

	n.availableCountries = append(n.availableCountries, n2)
}

func (n *Normalizer) Get(countryCode string) Normalize {
	for _, v := range n.availableCountries {
		if v.GetCountryCode() == countryCode {
			return v
		}
	}

	return n.fallback
}

func removeDuplicate(sliceList []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range sliceList {
		if item == "" {
			continue
		}
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
