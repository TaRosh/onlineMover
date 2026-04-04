package connection

import (
	"crypto/aes"
	"crypto/cipher"
)

func (conn *Conn) CreateEncryptor(sharedKey []byte) error {
	block, err := aes.NewCipher(sharedKey)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.gcm = gcm
	conn.state = Secure
	return nil
}
