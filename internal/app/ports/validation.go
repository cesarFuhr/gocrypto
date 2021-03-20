package ports

import (
	"time"

	"github.com/cesarFuhr/validator"
)

var (
	scopeV         = validator.NewStringValidator("scope", true, validator.StrLength(1, 50))
	expirationV    = validator.NewStringValidator("expiration", true, validator.StrDate(time.RFC3339))
	keyIDV         = validator.NewStringValidator("keyID", true, validator.StrUUID())
	dataV          = validator.NewStringValidator("data", true, validator.StrLength(1, 1000))
	encryptedDataV = validator.NewStringValidator("encryptedData", true, validator.StrLength(1, 4000))
)

type keysValidator struct{}

func (v keysValidator) PostValidator(ko keyOpts) error {
	if err := scopeV.Validate(ko.Scope); err != nil {
		return err
	}
	if err := expirationV.Validate(ko.Expiration); err != nil {
		return err
	}
	return nil
}

func (v keysValidator) GetValidator(keyID string) error {
	if err := keyIDV.Validate(keyID); err != nil {
		return err
	}
	return nil
}

type encryptValidator struct{}

func (v encryptValidator) PostValidator(eo encryptReqBody) error {
	if err := keyIDV.Validate(eo.KeyID); err != nil {
		return err
	}
	if err := dataV.Validate(eo.Data); err != nil {
		return err
	}
	return nil
}

type decryptValidator struct{}

func (v decryptValidator) PostValidator(do decryptReqBody) error {
	if err := keyIDV.Validate(do.KeyID); err != nil {
		return err
	}
	if err := encryptedDataV.Validate(do.EncryptedData); err != nil {
		return err
	}
	return nil
}
