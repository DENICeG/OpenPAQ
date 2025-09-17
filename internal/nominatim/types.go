package nominatim

import (
	"crypto/md5"
	"encoding/hex"
	"openPAQ/internal/normalization"
	"strings"
)

type NominatimResult struct {
	Address NominatimCoreResult `json:"address"`
}

type NominatimCoreResult struct {
	State            string `json:"state"`
	StateDistrict    string `json:"state_district"`
	Municipality     string `json:"municipality"`
	City             string `json:"city"`
	Town             string `json:"town"`
	Village          string `json:"village"`
	CityDistrict     string `json:"city_district"`
	District         string `json:"district"`
	Borough          string `json:"borough"`
	Suburb           string `json:"suburb"`
	Subdivision      string `json:"subdivision"`
	Hamlet           string `json:"hamlet"`
	Croft            string `json:"croft"`
	IsolatedDwelling string `json:"isolated_dwelling"`
	Neighbourhood    string `json:"neighbourhood"`
	Allotments       string `json:"allotments"`
	Quarter          string `json:"quarter"`
	Residential      string `json:"residential"`
	Farm             string `json:"farm"`
	Farmyard         string `json:"farmyard"`
	Industrial       string `json:"industrial"`
	Commercial       string `json:"commercial"`
	Retail           string `json:"retail"`
	Road             string `json:"road"`
	Park             string `json:"park"`
	Building         string `json:"building"`
	CityBlock        string `json:"city_block"`
	PostCode         string `json:"postcode"`
	CountryCode      string `json:"country_code"`
	County           string `json:"county"`
}

func (nr *NominatimCoreResult) parse() ParsedResult {
	p := ParsedResult{}
	if len(nr.State) > 0 {
		p.City = append(p.City, nr.State)
	}
	if len(nr.StateDistrict) > 0 {
		p.City = append(p.City, nr.StateDistrict)
	}
	if len(nr.Municipality) > 0 {
		p.City = append(p.City, nr.Municipality)
	}
	if len(nr.City) > 0 {
		p.City = append(p.City, nr.City)
	}
	if len(nr.Town) > 0 {
		p.City = append(p.City, nr.Town)
	}
	if len(nr.Village) > 0 {
		p.City = append(p.City, nr.Village)
	}
	if len(nr.CityDistrict) > 0 {
		p.City = append(p.City, nr.CityDistrict)
	}
	if len(nr.District) > 0 {
		p.City = append(p.City, nr.District)
	}
	if len(nr.Borough) > 0 {
		p.City = append(p.City, nr.Borough)
	}
	if len(nr.Suburb) > 0 {
		p.City = append(p.City, nr.Suburb)
	}
	if len(nr.Subdivision) > 0 {
		p.City = append(p.City, nr.Subdivision)
	}
	if len(nr.Hamlet) > 0 {
		p.City = append(p.City, nr.Hamlet)
		p.Street = append(p.Street, nr.Hamlet)
	}
	if len(nr.Croft) > 0 {
		p.City = append(p.City, nr.Croft)
	}
	if len(nr.IsolatedDwelling) > 0 {
		p.City = append(p.City, nr.IsolatedDwelling)
	}
	if len(nr.Neighbourhood) > 0 {
		p.City = append(p.City, nr.Neighbourhood)
	}
	if len(nr.Allotments) > 0 {
		p.City = append(p.City, nr.Allotments)
	}
	if len(nr.Quarter) > 0 {
		p.City = append(p.City, nr.Quarter)
	}
	if len(nr.Residential) > 0 {
		p.City = append(p.City, nr.Residential)
	}
	if len(nr.Farm) > 0 {
		p.City = append(p.City, nr.Farm)
	}
	if len(nr.Farmyard) > 0 {
		p.City = append(p.City, nr.Farmyard)
	}
	if len(nr.Industrial) > 0 {
		p.City = append(p.City, nr.Industrial)
	}
	if len(nr.Commercial) > 0 {
		p.City = append(p.City, nr.Commercial)
	}
	if len(nr.Retail) > 0 {
		p.City = append(p.City, nr.Retail)
	}
	if len(nr.County) > 0 {
		p.City = append(p.City, nr.County)
	}

	if len(nr.Road) > 0 {
		p.Street = append(p.Street, nr.Road)
	}
	if len(nr.Park) > 0 {
		p.Street = append(p.Street, nr.Park)
	}

	if len(nr.CityBlock) > 0 {
		p.Street = append(p.Street, nr.CityBlock)
	}
	if len(nr.Building) > 0 {
		p.Street = append(p.Street, nr.Building)
	}
	if len(nr.PostCode) > 0 {
		p.PostalCode = strings.ReplaceAll(strings.ToLower(nr.PostCode), " ", "")
	}
	if len(nr.CountryCode) > 0 {
		p.CountryCode = strings.ToLower(nr.CountryCode)
	}

	for i, v := range p.City {
		p.City[i] = strings.ToLower(v)
	}
	for i, v := range p.Street {
		p.Street[i] = strings.ToLower(v)
	}
	return p
}

type ParsedResult struct {
	Street      []string
	PostalCode  string
	City        []string
	CountryCode string
}

func (pr *ParsedResult) normalize(n normalization.Normalize) (ParsedResult, error) {
	ret := ParsedResult{}

	var err error
	ret.PostalCode, err = n.PostalCode(pr.PostalCode)
	if err != nil {
		return ret, err
	}

	for _, i := range pr.City {
		c, err := n.City(i)
		if err != nil {
			return ret, err
		}
		ret.City = append(ret.City, c)
	}

	for _, i := range pr.Street {
		s, err := n.Street(i)
		if err != nil {
			return ret, err
		}
		ret.Street = append(ret.Street, s...)
	}

	ret.CountryCode = pr.CountryCode

	return ret, nil
}

func (pr *ParsedResult) hash() string {
	res := pr.CountryCode
	res = res + pr.PostalCode
	res = res + strings.Join(pr.City, "")
	res = res + strings.Join(pr.Street, "")

	hash := md5.Sum([]byte(res))
	return hex.EncodeToString(hash[:])
}

type CityStreetResponse struct {
	City                  string
	Street                string
	WasPartialStreetMatch bool
	ParsedResult          ParsedResult
}
