package domain

type UserRepository interface {
	// Create writes passed user domain model to the users table.
	// A successful query execution sets user's id to the last inserted value, otherwise returns an error.
	// CreatedAt timestamp is set before execution.
	Create(user *User) error

	// Find searches for a certain user in the users table by its id field.
	// A successful search returns a non-nil user domain model. Otherwise, it returns an error.
	// If no user is found, it returns a nil instead of a user model and does not produce an error.
	Find(id int) (*User, error)

	// FindByEmail searches for a certain user in the users table by its id field.
	// A successful search returns a non-nil user domain model. Otherwise, it returns an error.
	// If no user is found, it returns a nil instead of a user model and does not produce an error.
	FindByEmail(email string) (*User, error)
}
