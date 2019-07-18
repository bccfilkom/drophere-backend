package stringgenerator

import (
	"strings"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/gofrs/uuid"
)

type myUUID struct{}

// NewUUID returns uuid token generator
func NewUUID() domain.StringGenerator {
	return &myUUID{}
}

// Generate generates random string
func (u *myUUID) Generate() string {
	token := uuid.Must(uuid.NewV4())
	return strings.ReplaceAll(token.String(), "-", "")
}
