package identity

type Repository interface {
	CreateUser(user *User) error
	GetUser(username string) (User, error)
}
