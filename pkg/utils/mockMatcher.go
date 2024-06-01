package utils

import "github.com/Pr3c10us/gbt/internals/domain/identity"

// Custom matcher for identity.User
func UserMatcher(expected *identity.User) func(*identity.User) bool {
	return func(actual *identity.User) bool {
		return actual.Username == expected.Username &&
			actual.SecurityQuestion == expected.SecurityQuestion
	}
}
