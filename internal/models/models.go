package models

import (
	"fmt"
	"sync"
	validator_module "validation_service/internal/validator"
	"validation_service/pkg/log"
)

type Agreement struct {
	service string
	logId   string
	fields  *sync.Map
	Errors  []string
	errors  chan string
}

func NewAgreement(service, logId string) *Agreement {
	a := &Agreement{
		service: service,
		logId:   logId,
		fields:  new(sync.Map),
		Errors:      make([]string, 0),
	}
	a.fields.Store("car", new(sync.Map))
	a.fields.Store("owner", new(sync.Map))
	a.fields.Store("insurer", new(sync.Map))
	a.fields.Store("agreement", new(sync.Map))

	return a
}

func (a *Agreement) Validate(data map[string]interface{}) {
	a.errors = make(chan string)
	waiter := new(sync.WaitGroup)

	if carFields, ok := data["car"]; ok {
		waiter.Add(1)
		go a.validateCarFields(carFields.(map[string]interface{}), waiter)
	}
	if ownerFields, ok := data["owner"]; ok {
		waiter.Add(1)
		go a.validateOwnerFields(ownerFields.(map[string]interface{}), waiter)
	}
	if insurerFields, ok := data["insurer"]; ok {
		waiter.Add(1)
		go a.validateInsurerFields(insurerFields.(map[string]interface{}), waiter)
	}
	if driversFieldsRaw, ok := data["drivers"]; ok {
		allDriversFields := driversFieldsRaw.([]interface{})

		for i, driversFields := range allDriversFields {
			waiter.Add(1)
			go a.validateDriversFields(i, driversFields.(map[string]interface{}), waiter)
		}
	}
	if agreementFields, ok := data["agreement"]; ok {
		waiter.Add(1)
		go a.validateAgreementFields(agreementFields.(map[string]interface{}), waiter)
	}

	errorWaiter := sync.WaitGroup{}
	errorWaiter.Add(1)
	go func() {
		for e := range a.errors {
			a.Errors = append(a.Errors, e)
		}

		if len(a.Errors) > 0 {
			log.Logger.Warnf(
				"%s/%s\n\nОшибки: %v\n\nПервоначальный запрос: %v", a.service, a.logId, a.Errors, data,
			)
		}

		errorWaiter.Done()
	}()

	waiter.Wait()
	close(a.errors)
	errorWaiter.Wait()
}

func (a *Agreement) validateCarFields(carFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	a.validate("car", "car", carFields)
}

func (a *Agreement) validateOwnerFields(ownerFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	a.validate("owner", "person", ownerFields)
}

func (a *Agreement) validateInsurerFields(insurerFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	a.validate("insurer", "person", insurerFields)
}

func (a *Agreement) validateDriversFields(i int, driversFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()

	key := fmt.Sprintf("driver_%d", i)

	a.fields.Store(key, new(sync.Map))
	a.validate(key, "driver", driversFields)
}

func (a *Agreement) validateAgreementFields(ownerFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	a.validate("agreement", "agreement", ownerFields)
}

func (a *Agreement) validate(object, validationClassPattern string, data map[string]interface{}) {
	validationPattern, err := validator_module.Registry.GetValidationPattern(validationClassPattern)
	if err != nil {
		log.Logger.Errorf("Нет шаблона под валидацию для: %s", validationClassPattern)
		return
	}

	validator, err := validator_module.Registry.GetValidator(validationPattern)
	if err != nil {
		log.Logger.Errorf("Ошибка при создании валидатора: %s", err.Error())
		return
	}

	validator.Init(object, a.fields, a.errors)

	validator.Validate(data)
}
