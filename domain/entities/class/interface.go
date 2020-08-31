package class

// IClass interface defines the methods that Class repository must implement
type IClass interface {
	Add(*Class) error
	Update(*Class) error
	Delete(int64) error
	Get(*Class) error
	GetAll(interface{}) error
}
