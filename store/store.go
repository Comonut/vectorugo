package store

//Store is the defintion for the storage implementation requirements
type Store interface {
	Get(id string) (Vector, error)
	Set(id string, vector Vector) error
	Delete(id string) error
	KNN(vector Vector, k int) (*[]Distance, error)
}
