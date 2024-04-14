package requests

type SignInRequest struct {
	Name string `json:"name"`
}

type SignInResponse struct {
	Token        string `json:"token,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}
