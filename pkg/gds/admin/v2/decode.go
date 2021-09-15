package admin

import (
	"encoding/json"
	"strings"
)

func (a *ResendAction) UnmarshalJSON(b []byte) (err error) {
	// Define a secondary type to avoid recursive call to UnmarshalJSON
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		panic(err)
	}

	// Ensure enum is white-space= and case-insensitive
	s = strings.TrimSpace(strings.ToLower(s))
	action := ResendAction(s)

	switch action {
	case ResendVerifyContact, ResendReview, ResendDeliverCerts, ResendRejection:
		*a = action
		return nil
	default:
		return ErrInvalidResendAction
	}
}
