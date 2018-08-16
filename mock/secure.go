package mock

// Secure mock
type Secure struct {
	PasswordFn    func(string, ...string) bool
	HashFn        func(string) string
	MatchesHashFn func(string, string) bool
}

// Password mock
func (s *Secure) Password(pw string, inputs ...string) bool {
	return s.PasswordFn(pw, inputs...)
}

// Hash mock
func (s *Secure) Hash(str string) string {
	return s.HashFn(str)
}

// MatchesHash mock
func (s *Secure) MatchesHash(hash, pw string) bool {
	return s.MatchesHashFn(hash, pw)
}
