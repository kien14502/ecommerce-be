package repo

type IUserRepository interface {
	FindOne(userID string) string
	FindAll() []string
	GetUserByEmail(email string) string
}

type userRepositoryType struct {
}

// FindAll implements [impl.IUserRepository].
func (u *userRepositoryType) FindAll() []string {
	panic("unimplemented")
}

// FindOne implements [impl.IUserRepository].
func (u *userRepositoryType) FindOne(userID string) string {
	panic("unimplemented")
}

// GetUserByEmail implements [impl.IUserRepository].
func (u *userRepositoryType) GetUserByEmail(email string) string {
	panic("unimplemented")
}

func NewUserRepository() IUserRepository {
	return &userRepositoryType{}
}
