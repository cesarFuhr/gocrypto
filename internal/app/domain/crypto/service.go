package crypto

// Service defines the methods of the crypto service
type Service interface {
	Encrypt(string, string) ([]byte, error)
	Decrypt(string, string) ([]byte, error)
}
