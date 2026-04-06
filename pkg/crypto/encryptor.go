package crypto

type Encryptor interface {
	Encrypt(plaintext string) (ciphertext string, err error)
	Decrypt(ciphertext string) (plaintext string, err error)
}
