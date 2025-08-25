package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func testmakeJWT(t *testing.T) {
	cases := []struct {
		inputUserID      uuid.UUID
		inputTokenString string
		inputExpiresAt   time.Duration
	}{
		{
			inputUserID:      uuid.MustParse("8c82c3d6-c292-4e0a-95c8-8476a04c12d2"),
			inputTokenString: "myToken",
			inputExpiresAt:   time.Duration(10 * time.Second),
		},
		{
			inputUserID:      uuid.MustParse("8c82c3d6-c292-4e0a-95c8-8476a04c12d2"),
			inputTokenString: "myToken",
			inputExpiresAt:   time.Duration(-10 * time.Second),
		},
	}

	for _, c := range cases {
		output, err := MakeJWT(c.inputUserID, c.inputTokenString)
		if err != nil {
			t.Error("Unable to create token")
		}

		userID, err := ValidateJWT(output, c.inputTokenString)
		if c.inputExpiresAt < 0 {
			if err == nil {
				t.Errorf("Expected error %v, but received nil response", jwt.ErrTokenExpired)
			} else if !errors.Is(err, jwt.ErrTokenExpired) {
				t.Errorf("Expected %v, but got %v", jwt.ErrTokenExpired, err)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Validation failed with error: %v", err)
			} else {
				if userID != c.inputUserID {
					t.Errorf("Expected userID: %v, but got %v", c.inputUserID, userID)
				}
			}
		}

		_, err = ValidateJWT(output, "badtokenstring")
		if err == nil {
			t.Errorf("Expected error %v, but received nil response", jwt.ErrSignatureInvalid)
		} else {
			if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				t.Errorf("Expected error %v, but got %v", jwt.ErrSignatureInvalid, err)
			}
		}

	}
}
