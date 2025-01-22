package models

// AddressValidationResponse represents the response for address validation
type AddressValidationResponse struct {
	Address          string `json:"address"`
	IsValid          bool   `json:"isValid"`
	HasValidChecksum bool   `json:"hasValidChecksum"`
	ChecksumAddress  string `json:"checksumAddress,omitempty"`
	Error            string `json:"error,omitempty"`
}
