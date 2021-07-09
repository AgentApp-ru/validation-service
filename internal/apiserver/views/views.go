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

func ValidateAgreement(bodyRaw []byte) ([]string, error) {
	var body map[string]interface{}
	if err := json.Unmarshal(bodyRaw, &body); err != nil {
		return []string{}, err
	}

	ps, agreementID := getPsAndAgreementID(body)
	agreement := models.NewAgreement(ps, agreementID)
	agreement.Validate(body)

	return agreement.Errors, nil
}

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
