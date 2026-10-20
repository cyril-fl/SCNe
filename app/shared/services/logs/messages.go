package logs

import (
	"errors"
	"fmt"
	"solopg/app/shared/services/i19n"
)

// Success
func Success(id string, data ...map[string]any) string {
	SilentSuccess(id, data...)
	return CeaseSuccess(id, data...)
}

func SilentSuccess(id string, data ...map[string]any) {
	registerFromTemplate(Template{
		Type:    INFO,
		Message: i19n.Unlocalize(id, data...),
	})
}

func CeaseSuccess(id string, data ...map[string]any) string {
	return i19n.Localize(id, data...)
}

// Error
func Error(id string, data ...map[string]any) error {
	SilentError(id, data...)
	return CeaseError(id, data...)
}

func MildError(id string, data ...map[string]any) error {
	SilentWarning(id, data...)
	return CeaseError(id, data...)
}

func SoftError(id string, data ...map[string]any) error {
	SilentInfo(id, data...)
	return CeaseError(id, data...)
}

func SilentError(id string, data ...map[string]any) {
	registerFromTemplate(Template{
		Type:    ERR,
		Message: i19n.Unlocalize(id, data...),
	})
}

func RawError(format string, a ...any) error {
	err := fmt.Errorf(format, a...)

	registerFromTemplate(Template{
		Type:    ERR,
		Message: err.Error(),
	})

	return err
}

func CeaseError(id string, data ...map[string]any) error {
	return errors.New(i19n.Localize(id, data...))
}

func CatalogError(id string, data ...map[string]any) error {
	registerFromTemplate(Template{
		Type:    ERR,
		Message: i19n.Unlocalize(id, data...),
	})

	return errors.New(i19n.Unlocalize(id, data...))
}

// Info
func Info(id string, data ...map[string]any) string {
	SilentInfo(id, data...)
	return i19n.Localize(id, data...)
}

func SilentInfo(id string, data ...map[string]any) {
	registerFromTemplate(Template{
		Type:    INFO,
		Message: i19n.Unlocalize(id, data...),
	})
}

// Warning
func Warning(id string, data ...map[string]any) string {
	SilentWarning(id, data...)
	return i19n.Localize(id, data...)
}

func SilentWarning(id string, data ...map[string]any) {
	registerFromTemplate(Template{
		Type:    WARN,
		Message: i19n.Unlocalize(id, data...),
	})
}
