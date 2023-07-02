package handlers

import "fmt"

type MissingRequestParamError struct {
	ParamName string
}

func NewMissingRequestParamError(paramName string) error {
	return MissingRequestParamError{
		ParamName: paramName,
	}
}

func (e MissingRequestParamError) Error() string {
	return fmt.Sprintf("U requestu nije pronađen parametar %s", e.ParamName)
}

type WrongParamFormatError struct {
	ParamName  string
	ParamValue string
}

func NewWrongParamFormatErrorError(paramName, paramValue string) error {
	return WrongParamFormatError{
		ParamName:  paramName,
		ParamValue: paramValue,
	}
}

func (e WrongParamFormatError) Error() string {
	return fmt.Sprintf("Parametar %s (%s) nije u očekivanom formatu", e.ParamName, e.ParamValue)
}

type JSONDecodeError struct {
	StructName string
}

func NewJSONDecodeError(structName string) error {
	return JSONDecodeError{
		StructName: structName,
	}
}

func (e JSONDecodeError) Error() string {
	return fmt.Sprintf("Greška prilikom dekodovanja strukture %s", e.StructName)
}
