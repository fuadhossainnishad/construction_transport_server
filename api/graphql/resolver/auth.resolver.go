package resolver

type AuthResolver struct{}

func (r *AuthResolver) Login() string {
	return "jwt_token"
}
