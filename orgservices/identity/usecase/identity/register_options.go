package identity

import "github.com/timoth-y/chainmetric-network/orgservices/identity/model"

type (
	// RegistrationOption allows passing parameters for Register method.
	RegistrationOption interface {
		Apply(*model.User)
	}

	// RegisterOptionFunc is a function that mutates model during Register execution.
	RegisterOptionFunc func(*model.User)
)

// Apply calls RegisterOptionFunc on model.
func (f RegisterOptionFunc) Apply(user *model.User) {
	f(user)
}

// WithName creates user with given `firstname` and `lastname`.
func WithName(firstname, lastname string) RegistrationOption {
	return RegisterOptionFunc(func(u *model.User) {
		u.Firstname = firstname
		u.Lastname = lastname
	})
}

// WithEmail creates user with given `email`.
func WithEmail(email string) RegistrationOption {
	return RegisterOptionFunc(func(u *model.User) {
		u.Email = email
	})
}
