package models

import (
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// KeysTable is the name of the table the [Keys] model is stored in.
const KeysTable = "keys"

// Keys represents a distillery ssh key
type Keys struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	User string `gorm:"column:user;not null"` // username of the ssh key

	Signature []byte `gorm:"column:signature;not null"` // signature of the ssh key
	Comment   string `gorm:"column:comment"`
}

// PublicKey returns the public key corresponding to this keys.
// If the key cannot be parsed, returns nil.
func (keys *Keys) PublicKey() ssh.PublicKey {
	key, err := ssh.ParsePublicKey(keys.Signature)
	if err != nil {
		return nil
	}
	return key
}

func (keys *Keys) SignatureString() string {
	// try to get the public key
	key := keys.PublicKey()
	if key == nil {
		return ""
	}

	// marshal the key!
	return string(gossh.MarshalAuthorizedKey(key))
}

// SetPublicKey stores a specific public key in this key
func (keys *Keys) SetPublicKey(key ssh.PublicKey) {
	keys.Signature = key.Marshal()
}
