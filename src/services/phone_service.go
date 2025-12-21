package services

import (
	"api-i18n/main/src/dto/responses"
	"sort"
	"strconv"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// GetTerritoryPhoneCodes retrieves phone codes for all supported territories.
func GetTerritoryPhoneCodes(localeID string) (*[]responses.PhoneCodeLookup, error) {
	phoneCodes := make([]responses.PhoneCodeLookup, 0)

	regions := phonenumbers.GetSupportedRegions()
	territories, err := GetTerritoriesLookup(localeID, nil, nil)
	if err != nil {
		return nil, err
	}

	territoriesMap := make(map[string]string)
	for _, territory := range *territories {
		territoriesMap[territory.TerritoryID] = territory.Name
	}

	for region := range regions {
		code := phonenumbers.GetCountryCodeForRegion(region)
		if code <= 0 {
			continue
		}

		territoryName := region
		if name, exists := territoriesMap[region]; exists {
			territoryName = name
		}

		phoneCode := responses.PhoneCodeLookup{
			Code:      "+" + strconv.Itoa(code),
			Territory: region,
			Name:      territoryName,
		}
		phoneCodes = append(phoneCodes, phoneCode)
	}

	// Sort by numeric dial code, then by name, then by territory
	sort.SliceStable(phoneCodes, func(i, j int) bool {
		ci, _ := strconv.Atoi(strings.TrimPrefix(phoneCodes[i].Code, "+"))
		cj, _ := strconv.Atoi(strings.TrimPrefix(phoneCodes[j].Code, "+"))
		if ci == cj {
			if phoneCodes[i].Name == phoneCodes[j].Name {
				return phoneCodes[i].Territory < phoneCodes[j].Territory
			}
			return phoneCodes[i].Name < phoneCodes[j].Name
		}
		return ci < cj
	})

	return &phoneCodes, nil
}

// ValidatePhoneNumber checks if the given phone number is valid for the specified region.
func ValidatePhoneNumber(number string, region *string) (bool, error) {
	n, r := preformatNumberAndRegion(number, region)
	if n == "" {
		return false, nil
	}

	// TODO: The phonenumbers library does not support to disable region checking. For now region is required from caller.
	parsedNumber, err := phonenumbers.Parse(n, r)
	if err != nil {
		return false, err
	}

	// If region is provided, validate for that region; otherwise validate globally.
	var isValid bool
	if r != "" {
		isValid = phonenumbers.IsValidNumberForRegion(parsedNumber, r)
	} else {
		isValid = phonenumbers.IsValidNumber(parsedNumber)
	}

	return isValid, nil
}

// FormatPhoneNumber formats the given phone number into various formats and provides validation info.
func FormatPhoneNumber(number string, region *string) (*responses.PhoneNumberFormat, error) {
	n, r := preformatNumberAndRegion(number, region)
	if n == "" {
		return nil, nil
	}

	// TODO: The phonenumbers library does not support to disable region checking. For now region is required from caller.
	parsedNumber, err := phonenumbers.Parse(n, r)
	if err != nil {
		return nil, err
	}

	var regionCode, rfc3966, e164, national, international, numberType *string
	var countryCode *int

	// If region is provided, validate for that region; otherwise validate globally.
	var isValid bool
	if r != "" {
		isValid = phonenumbers.IsValidNumberForRegion(parsedNumber, r)
	} else {
		isValid = phonenumbers.IsValidNumber(parsedNumber)
	}
	isPossible := phonenumbers.IsPossibleNumber(parsedNumber)
	regCode := phonenumbers.GetRegionCodeForNumber(parsedNumber)

	// Only set region if it's a specific region (exclude unknown "ZZ").
	if regCode != "" && regCode != "ZZ" {
		regionCode = &regCode
		ctyCode := phonenumbers.GetCountryCodeForRegion(regCode)
		if ctyCode > 0 {
			countryCode = &ctyCode
		}
	}

	nType := phonenumbers.GetNumberType(parsedNumber)
	ts := phoneNumberTypeToString(nType)
	if ts != "" {
		numberType = &ts
	}

	if isValid {
		result := phonenumbers.Format(parsedNumber, phonenumbers.E164)
		e164 = &result

		result = phonenumbers.Format(parsedNumber, phonenumbers.RFC3966)
		rfc3966 = &result
	}

	if isPossible {
		result := phonenumbers.Format(parsedNumber, phonenumbers.NATIONAL)
		national = &result

		result = phonenumbers.Format(parsedNumber, phonenumbers.INTERNATIONAL)
		international = &result
	}

	phoneNumberFormat := &responses.PhoneNumberFormat{
		IsValid:       isValid,
		IsPossible:    isPossible,
		RFC3966:       rfc3966,
		E164:          e164,
		National:      national,
		International: international,
		Region:        regionCode,
		Code:          countryCode,
		Type:          numberType,
	}

	return phoneNumberFormat, nil
}

// preformatNumberAndRegion trims spaces and ensures the number starts with '+' if no region is provided.
func preformatNumberAndRegion(number string, region *string) (string, string) {
	n := strings.TrimSpace(number)
	r := ""

	if region != nil {
		// Normalize region to uppercase ISO 3166-1 alpha-2 when provided
		r = strings.ToUpper(strings.TrimSpace(*region))
	}
	if r == "" {
		// If no region is given, require international format; prefix '+' when
		// we have a non-empty number that doesn't already have it.
		if n == "" {
			return "", r
		}
		if n[0] != '+' {
			n = "+" + n
		}
	}

	return n, r
}

func phoneNumberTypeToString(t phonenumbers.PhoneNumberType) string {
	switch t {
	case phonenumbers.FIXED_LINE:
		return "Fixed line"
	case phonenumbers.MOBILE:
		return "Mobile"
	case phonenumbers.FIXED_LINE_OR_MOBILE:
		return "Fixed line or mobile"
	case phonenumbers.TOLL_FREE:
		return "Toll free"
	case phonenumbers.PREMIUM_RATE:
		return "Premium rate"
	case phonenumbers.SHARED_COST:
		return "Shared cost"
	case phonenumbers.VOIP:
		return "VoIP"
	case phonenumbers.PERSONAL_NUMBER:
		return "Personal number"
	case phonenumbers.PAGER:
		return "Pager"
	case phonenumbers.UAN:
		return "UAN"
	case phonenumbers.VOICEMAIL:
		return "Voicemail"
	case phonenumbers.UNKNOWN:
		return ""
	default:
		return ""
	}
}
