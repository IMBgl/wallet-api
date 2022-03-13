package validator

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

const TYPE_FLOAT = "float"
const TYPE_STRING = "string"
const TYPE_UUID = "uuid"
const DEFAULT_MESSAGE = "validation error"

var DefaultMessages map[string]string = map[string]string{
	"float":  "%s value must be float",
	"string": "%s value must be string",
	"uuid":   "%s value must be uuid",
}

func GetMessage(validationType string) string {
	if m, ok := DefaultMessages[validationType]; ok {
		return m
	}
	return DEFAULT_MESSAGE
}

func GetFieldMessage(validationType, fieldName string) string {
	msg := GetMessage(validationType)
	if msg != DEFAULT_MESSAGE {
		msg = fmt.Sprintf(msg, fieldName)
	}

	return msg
}

func Float32(value interface{}, name string) (float32, error) {
	switch t := value.(type) {
	case string:
		floatVal, err := strconv.ParseFloat(t, 32)
		if err != nil {
			return 0.0, errors.New(GetFieldMessage(TYPE_FLOAT, name))
		}
		return float32(floatVal), nil
	case float32:
		return t, nil
	case float64:
		return float32(t), nil
	}

	return 0.0, errors.New(fmt.Sprintf("%s value must be float", name))
}

func Uuid(value interface{}, name string) (uuid.UUID, error) {
	switch t := value.(type) {
	case string:
		uuidVal, err := uuid.Parse(t)
		if err != nil {
			return uuid.New(), errors.New(GetFieldMessage(TYPE_UUID, name))
		}

		return uuidVal, nil
	}

	return uuid.New(), errors.New(GetFieldMessage(TYPE_UUID, name))
}

func String(value interface{}, name string) (string, error) {
	switch t := value.(type) {
	case string:
		return t, nil
	}

	return "", errors.New(GetFieldMessage(TYPE_STRING, name))
}
