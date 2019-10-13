package security

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/twinj/uuid"

)

func TokenHash(text string) string {

	hasher := md5.New()
	hasher.Write([]byte(text))
	theHash :=  hex.EncodeToString(hasher.Sum(nil))

	//also use uuid
	u := uuid.NewV4()
	theToken := theHash + u.String()

	return theToken
}