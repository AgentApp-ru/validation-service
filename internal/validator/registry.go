package validator

import (
	"validation_service/pkg/storage"
)

type registry struct {
	storage storage.Storage
}

var Registry *registry

func Init(store storage.Storage) {
	Registry = &registry{
		storage: store,
	}
}

func (r *registry) GetValidationPattern(object string) ([]byte, error) {
	var (
		rawData []byte
		err     error
	)

	rawData, err = r.storage.Get(object)
	return rawData, err
}

func (r *registry) GetValidator(data []byte) (*validator, error) {
	vc := new(validator)

	if err := vc.New(data); err != nil {
		return nil, err
	}

	return vc, nil
}
