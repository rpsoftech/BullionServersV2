package interfaces

type TokenResponseBody struct {
	AccessToken   string `json:"accessToken"`
	RefreshToken  string `json:"refreshToken"`
	FirebaseToken string `json:"firebaseToken"`
}
