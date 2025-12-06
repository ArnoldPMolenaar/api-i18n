package utils

import (
	"api-i18n/main/src/services"
	"strings"
)

// ResolveLocaleId attempts to resolve a localeID by checking availability and
// progressively stripping trailing segments separated by '-' until a match is found.
// Examples:
//
//	en -> en
//	en-ER -> en-ER, fallback to en if en-ER not available
//	en-polyton -> en-polyton, fallback to en if en-polyton not available
//	az-Arab-IQ -> az-Arab-IQ, fallback to az-Arab, then az
//
// Returns a pointer to the resolved locale ID or nil if none are available.
func ResolveLocaleId(localeID string) *string {
	if localeID == "" {
		return nil
	}

	// First, try the provided localeID as-is.
	if ok, err := services.IsLocaleAvailable(localeID); err == nil && ok {
		return &localeID
	}

	// Split into segments and iteratively remove the last segment.
	parts := strings.Split(localeID, "-")
	if len(parts) == 0 {
		return nil
	}

	// Iteratively try with fewer trailing segments: language[-script[-territory|variant]]
	for i := len(parts) - 1; i >= 1; i-- {
		candidate := strings.Join(parts[:i], "-")
		if ok, err := services.IsLocaleAvailable(candidate); err == nil && ok {
			// Return a pointer to a new string variable to avoid referencing the loop variable.
			resolved := candidate
			return &resolved
		}
	}

	// No match found.
	return nil
}
