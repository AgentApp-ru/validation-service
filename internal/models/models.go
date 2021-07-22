package models

import (
	"fmt"
	"sync"
	validator_module "validation_service/internal/validator"
	"validation_service/pkg/log"
)

type Agreement struct {
	service      string
	logId        string
	fields       *sync.Map
	Errors       []string
	errors       chan string
	AbsentFields []string
	absentFields chan string
}

func NewAgreement(service, logId string) *Agreement {
	a := &Agreement{
		service:      service,
		logId:        logId,
		fields:       new(sync.Map),
		Errors:       make([]string, 0),
		AbsentFields: make([]string, 0),
	}
	a.fields.Store("car", new(sync.Map))
	a.fields.Store("owner", new(sync.Map))
	a.fields.Store("insurer", new(sync.Map))
	a.fields.Store("agreement", new(sync.Map))

	return a
}

func (a *Agreement) Validate(data map[string]interface{}) {
	a.errors = make(chan string)
	a.absentFields = make(chan string)
	waiter := new(sync.WaitGroup)

	carFields := data["car"]
	waiter.Add(1)
	go a.validateFields("car", carFields, waiter)

	ownerFields := data["owner"]
	waiter.Add(1)
	go a.validateFields("owner", ownerFields, waiter)

	insurerFields := data["insurer"]
	waiter.Add(1)
	go a.validateFields("insurer", insurerFields, waiter)

	agreementFields := data["agreement"]
	waiter.Add(1)
	go a.validateFields("agreement", agreementFields, waiter)

	if driversFieldsRaw, ok := data["drivers"]; ok {
		allDriversFields := driversFieldsRaw.([]interface{})

		for i, driversFields := range allDriversFields {
			waiter.Add(1)
			go a.validateDriversFields(i, driversFields.(map[string]interface{}), waiter)
		}
	}

	absentWaiter := sync.WaitGroup{}
	absentWaiter.Add(1)
	go func() {
		for e := range a.absentFields {
			a.AbsentFields = append(a.AbsentFields, e)
		}

		absentWaiter.Done()
	}()

	errorWaiter := sync.WaitGroup{}
	errorWaiter.Add(1)
	go func() {
		for e := range a.errors {
			a.Errors = append(a.Errors, e)
		}

		errorWaiter.Done()
	}()

	waiter.Wait()
	close(a.absentFields)
	close(a.errors)
	absentWaiter.Wait()
	errorWaiter.Wait()

	if len(a.Errors) > 0 || len(a.AbsentFields) > 0 {
		log.Logger.Warnf(
			"%s/%s\n\nОтсутствующие поля: %v\n\nОшибки: %v\n\nПервоначальный запрос: %v", a.service, a.logId, a.AbsentFields, a.Errors, data,
		)
	}
}

func (a *Agreement) validateFields(object string, rawFields interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()

	var fields map[string]interface{}
	if rawFields == nil {
		fields = make(map[string]interface{})
	} else {
		fields = rawFields.(map[string]interface{})
	}

	a.validate(object, object, fields)
}

func (a *Agreement) validateDriversFields(i int, driversFields map[string]interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()

	key := fmt.Sprintf("driver_%d", i)

	a.fields.Store(key, new(sync.Map))
	a.validate(key, "driver", driversFields)
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

	validator.Init(object, a.fields, a.errors, a.absentFields)

	validator.Validate(data)
}
