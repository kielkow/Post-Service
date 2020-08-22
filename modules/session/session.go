package session

// Session struct
type Session struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Token struct
type Token struct {
	Token string `json:"token"`
}
