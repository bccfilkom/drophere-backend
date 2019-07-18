package stringgenerator

import "github.com/bccfilkom/drophere-go/domain"

var preset string

type mockStringGenerator struct{}

// SetMockResult set the string that Generate returns
func SetMockResult(s string) {
	preset = s
}

// NewMock returns new mockStringGenerator
func NewMock() domain.StringGenerator {
	return &mockStringGenerator{}
}

// Generate returns pre-set string
func (m *mockStringGenerator) Generate() string {
	return preset
}
