package user

// IUser interface defines the methods that User repository must implement
type IUser interface {
	Add(*User) error
	Update(*User) error
	Delete(int64) error
	Get(*User) error
	GetAll(interface{}) error
}
