package storage


type Storage interface {
	Get(string) ([]byte, error)
}