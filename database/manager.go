package database

type Manager interface {
	Retrieve(id int) (interface{}, error)
	Delete(id int) (interface{}, error)
	Find() ([]interface{}, error)
	Create(interface{}) (interface{}, error)
	Update(interface{}) (interface{}, error)
}
