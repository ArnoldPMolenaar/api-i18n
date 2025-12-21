package responses

type PhoneCodeLookup struct {
	Code      string `json:"code"`
	Territory string `json:"territory"`
	Name      string `json:"name"`
}

// SetPhoneCodeLookup sets the phone code lookup details.
func (pcl *PhoneCodeLookup) SetPhoneCodeLookup(code, territory, name string) {
	pcl.Code = code
	pcl.Territory = territory
	pcl.Name = name
}
