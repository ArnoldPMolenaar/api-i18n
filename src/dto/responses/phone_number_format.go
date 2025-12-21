package responses

type PhoneNumberFormat struct {
	IsValid       bool    `json:"isValid"`
	IsPossible    bool    `json:"isPossible"`
	RFC3966       *string `json:"rfc3966"`
	E164          *string `json:"e164"`
	National      *string `json:"national"`
	International *string `json:"international"`
	Region        *string `json:"region"`
	Code          *int    `json:"code"`
	Type          *string `json:"type"`
}
