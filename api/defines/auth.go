package defines

type Auth struct {
	EncryptedString string `json:"encrypted_string" form:"encrypted_string"`
	ClientID        string `json:"client_id" form:"client_id"`
}
