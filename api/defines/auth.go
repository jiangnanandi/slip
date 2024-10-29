package types

type Auth struct {
	EncryptedString string `json:"encrypted_string"`
	ClientID        string `json:"client_id"`
}
