package models

import (
	"sync"
	"time"
	validator_module "validation_service/internal/validator"
	"validation_service/pkg/log"
)

type Agreement struct {
	ps string
	agreementID string
	fields *sync.Map
	Errors []string
	errors chan string
}

func NewAgreement(ps, agreementID string) *Agreement {
	a := &Agreement{
		ps: ps,
		agreementID: agreementID,
		fields: new(sync.Map),
		Errors: make([]string, 0),
	}
	a.fields.Store("car", new(sync.Map))
	a.fields.Store("owner", new(sync.Map))
	a.fields.Store("insurer", new(sync.Map))
	a.fields.Store("drivers", new(sync.Map))
	a.fields.Store("agreement", new(sync.Map))

	return a
}

func (a *Agreement) Validate(data map[string]map[string]interface{}) {
	a.errors = make(chan string)
	defer close(a.errors)

	waiter := new(sync.WaitGroup)

	if carFields, ok := data["car"]; ok {
		waiter.Add(1)
		go a.validateCarFields(carFields, waiter)
	}
	if ownerFields, ok := data["owner"]; ok {
		waiter.Add(1)
		go a.validateOwnerFields(ownerFields, waiter)
	}
	if insurerFields, ok := data["insurer"]; ok {
		waiter.Add(1)
		go a.validateInsurerFields(insurerFields, waiter)
	}
	// if driversFields, ok := data["drivers"]; ok {
	// 	waiter.Add(1)
	// 	go a.validateDriversFields(driversFields)
	// }
	if generalFields, ok := data["general"]; ok {
		waiter.Add(1)
		go a.validateGeneralFields(generalFields, waiter)
	}

	go func() {
		for e := range a.errors {
			log.Logger.Warnf("%s/%s. Не прошла валидация %s", a.ps, a.agreementID, e)
			a.Errors = append(a.Errors, e)
		}
	}()

	waiter.Wait()
	time.Sleep(2 * time.Second) // Дожидаемся, что все ошибки созранятся
}

func (a *Agreement) validateCarFields(carFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	a.validate("car", "car", carFields)
}

func (a *Agreement) validateOwnerFields(ownerFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	a.validate("owner", "person", ownerFields)
}

func (a *Agreement) validateInsurerFields(ownerFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	a.validate("insurer", "person", ownerFields)
}

func (a *Agreement) validateGeneralFields(ownerFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	a.validate("agreement", "general", ownerFields)
}

func (a *Agreement) validate(object, validationClassPattern string, data map[string]interface{}) {
	rawValidator, err := validator_module.Validator.GetRaw(validationClassPattern)
	if err != nil {
		log.Logger.Errorf("Нет валидаторов для: %s", validationClassPattern)
		return
	}

	validatorClass, err := validator_module.Validator.GetValidatorClass(rawValidator)
	if err != nil {
		log.Logger.Errorf("Ошибка при создании класса валидатора: %s", err.Error())
		return
	}

	validatorClass.Validate(object, a.fields, data, a.errors)
}
