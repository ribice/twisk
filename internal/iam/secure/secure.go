package secure

import (
	zxcvbn "github.com/nbutton23/zxcvbn-go"
	"golang.org/x/crypto/bcrypt"
)

// New creates new password service
func New(minStr int) *Service {
	return &Service{minStr: minStr}
}

// Service contains password related methods
type Service struct {
	minStr   int
	appWords []string
}

// Password checks whether password is secure enough using zxcvbn library
func (s *Service) Password(pw string, inputs ...string) bool {
	inputs = append(inputs, s.appWords...)
	pwStrength := zxcvbn.PasswordStrength(pw, inputs)
	return !(pwStrength.Score < s.minStr)
}

// Hash hashes the password using bcrypt
func (s *Service) Hash(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// MatchesHash checks whether hash matches with password. Returns true if hash and password match.
func (s *Service) MatchesHash(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
