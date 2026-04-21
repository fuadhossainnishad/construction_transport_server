package redis

type SessionStore struct{}

func (s *SessionStore) Set(userId int64, token string) error {
	return nil
}

func (s *SessionStore) Get(userId int64) (string, error) {
	return "", nil
}
