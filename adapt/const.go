package adapt

import (
	"errors"
	"reflect"
)

type tagName string
type tagValue string

type tagFunction func(tagValue, reflect.Value) error

type tagsList map[tagName]tagValue

const (
	SET_DELIMITER = "||"
	VAL_DELIMITER = "**"

	TAG_VALUE = "value"
	TAG_JSON  = "json"
	TAG_INFO  = "info"

	RST_MIN       = "rst-min"
	RST_MAX       = "rst-max"
	RST_REGEX     = "rst-regex"
	RST_DEFAULT   = "rst-default"
	RST_CHOICE    = "rst-choice"
	RST_FORBIDDEN = "rst-forbidden"

	// removed unused VLD_* constants
)

var (
	ErrNotStruct   = errors.New("argument is not a struct")
	ErrInvalidTags = errors.New("invalid struct tags")
)

var tagsMap = map[tagName]tagFunction{
	RST_MIN:       adaptMin,
	RST_MAX:       adaptMax,
	RST_REGEX:     adaptRegex,
	RST_DEFAULT:   adaptDefault,
	RST_CHOICE:    adaptChoice,
	RST_FORBIDDEN: adaptForbidden,
}

// В последние 3 тега добавить разделители для строковых значений
// Разобраться с float32
// Добавить в тесты пустые поля, проверить на конфликт тегов
