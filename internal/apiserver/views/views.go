package views

import (
	"encoding/json"
	"validation_service/internal/models"
	"validation_service/internal/validator"
)

func Ping() string {
	return "pong"
}

func GetValidationPattern(object string) (interface{}, error) {
	var (
		result  interface{}
		rawData []byte
		err     error
	)

	rawData, err = validator.Registry.GetValidationPattern(object)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawData, &result)

	return result, err
}

func ValidateAgreement(bodyRaw []byte) ([]string, []string, error) {
	var (
		body    map[string]interface{}
		service string
		logId   string
	)
	if err := json.Unmarshal(bodyRaw, &body); err != nil {
		return []string{}, []string{}, err
	}

	ps, agreementID := getPsAndAgreementID(body)
	if ps == "" && agreementID == "" {
		service, logId = getServiceAndLogID(body)
	} else { // deprecated branch
		service = ps
		logId = agreementID
	}

	agreement := models.NewAgreement(service, logId)
	agreement.Validate(body)

	return agreement.AbsentFields, agreement.Errors, nil
}

func getServiceAndLogID(b map[string]interface{}) (string, string) {
	var service, logId string

	generalLoaded, ok := b["general"]
	if ok {
		general := generalLoaded.(map[string]interface{})
		serviceRaw, ok := general["service"]
		if ok {
			service = serviceRaw.(string)
		}
		logIRaw, ok := general["log_id"]
		if ok {
			logId = logIRaw.(string)
		}

	}

	return service, logId
}

// deprecated
func getPsAndAgreementID(b map[string]interface{}) (string, string) {
	var ps, agreementID string
	generalLoaded, ok := b["general"]
	if ok {
		general := generalLoaded.(map[string]interface{})
		psRaw, ok := general["ps"]
		if ok {
			ps = psRaw.(string)
		}
		agreementIdRaw, ok := general["agreement_id"]
		if ok {
			agreementID = agreementIdRaw.(string)
		}
	}

	return ps, agreementID
}
