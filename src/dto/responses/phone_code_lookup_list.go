package responses

type PhoneCodeLookupList struct {
	Codes []PhoneCodeLookup `json:"codes"`
}

// SetPhoneCodeLookupList sets the list of phone code lookups.
func (pcll *PhoneCodeLookupList) SetPhoneCodeLookupList(codes *[]PhoneCodeLookup) {
	pcll.Codes = *codes
}
