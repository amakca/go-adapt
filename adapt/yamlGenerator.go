package adapter

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// GenerateStructYAML генерирует YAML файл структуры с комментариями из структурных тегов
// Функция получает на вход структуру и возвращает строку с YAML представлением
func GenerateStructYAML(input any) (string, error) {
	inputValue := reflect.ValueOf(input)

	if reflect.Indirect(inputValue).Kind() != reflect.Struct {
		return "", ErrNotStruct
	}

	var result strings.Builder
	result.WriteString("# Generated YAML structure with RST tags comments\n\n")

	if err := generateStructYAMLRecursive(inputValue, "", &result, 0); err != nil {
		return "", err
	}

	return result.String(), nil
}

// GenerateStructYAMLFile генерирует YAML файл структуры с комментариями из структурных тегов
// Функция получает на вход структуру и имя файла, создает .yaml файл
func GenerateStructYAMLFile(input any, filename string) error {
	yaml, err := GenerateStructYAML(input)
	if err != nil {
		return err
	}

	// Добавляем расширение .yaml если его нет
	if !strings.HasSuffix(filename, ".yaml") && !strings.HasSuffix(filename, ".yml") {
		filename += ".yaml"
	}

	// Записываем в файл
	err = os.WriteFile(filename, []byte(yaml), 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл %s: %w", filename, err)
	}

	return nil
}

// generateStructYAMLRecursive рекурсивно генерирует YAML для структуры
func generateStructYAMLRecursive(input reflect.Value, fieldName string, result *strings.Builder, indent int) error {
	// Обрабатываем указатели
	if input.Kind() == reflect.Ptr {
		if input.IsNil() {
			// Для nil указателей генерируем null
			indentStr := strings.Repeat("  ", indent)
			result.WriteString(fmt.Sprintf("%s%s: null\n", indentStr, fieldName))
			return nil
		}
		input = input.Elem()
	}

	// Обрабатываем интерфейсы
	if input.Kind() == reflect.Interface {
		if input.IsNil() {
			indentStr := strings.Repeat("  ", indent)
			result.WriteString(fmt.Sprintf("%s%s: null\n", indentStr, fieldName))
			return nil
		}
		input = input.Elem()
	}

	indentStr := strings.Repeat("  ", indent)

	switch input.Kind() {
	case reflect.Struct:
		if fieldName != "" {
			result.WriteString(fmt.Sprintf("%s%s:\n", indentStr, fieldName))
		}

		inputType := input.Type()
		for i := 0; i < input.NumField(); i++ {
			field := inputType.Field(i)
			value := input.Field(i)

			// Пропускаем неэкспортируемые поля
			if !field.IsExported() {
				continue
			}

			// Получаем имя поля из json тега или используем имя поля (по умолчанию в нижнем регистре)
			jsonTag := field.Tag.Get("json")
			name := strings.ToLower(field.Name)
			if jsonTag != "" && jsonTag != "-" {
				// Убираем omitempty если есть
				if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
					name = jsonTag[:commaIdx]
				} else {
					name = jsonTag
				}
			}

			// Комментарии печатаем без отступа
			comment := generateCommentFromTags(field.Tag)
			if comment != "" {
				result.WriteString(fmt.Sprintf("# %s\n", comment))
			}

			if err := generateStructYAMLRecursive(value, name, result, indent+1); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		if fieldName != "" {
			result.WriteString(fmt.Sprintf("%s%s:\n", indentStr, fieldName))
		}

		indentNext := strings.Repeat("  ", indent+1)
		for i := 0; i < input.Len(); i++ {
			val := input.Index(i)
			// Handle pointer/interface elements
			if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
				if val.IsNil() {
					result.WriteString(fmt.Sprintf("%s- null\n", indentNext))
					continue
				}
				val = val.Elem()
			}

			if isSimpleKind(val.Kind()) {
				result.WriteString(fmt.Sprintf("%s- %s\n", indentNext, formatValue(val)))
			} else {
				result.WriteString(fmt.Sprintf("%s- \n", indentNext))
				if err := generateStructYAMLRecursive(val, "", result, indent+2); err != nil {
					return err
				}
			}
		}

	case reflect.Map:
		if fieldName != "" {
			result.WriteString(fmt.Sprintf("%s%s:\n", indentStr, fieldName))
		}

		// Sort keys for stable output
		keys := input.MapKeys()
		keyStrs := make([]string, len(keys))
		for i, k := range keys {
			keyStrs[i] = fmt.Sprintf("%v", k.Interface())
		}
		sort.Strings(keyStrs)
		for _, keyStr := range keyStrs {
			// Try direct string key first
			val := input.MapIndex(reflect.ValueOf(keyStr))
			if !val.IsValid() {
				// Fallback for non-string keys: match by stringified form
				for _, k := range keys {
					if fmt.Sprintf("%v", k.Interface()) == keyStr {
						val = input.MapIndex(k)
						break
					}
				}
			}
			if err := generateStructYAMLRecursive(val, keyStr, result, indent+1); err != nil {
				return err
			}
		}

	default:
		// Простые типы
		if fieldName != "" {
			valueStr := formatValue(input)
			result.WriteString(fmt.Sprintf("%s%s: %s\n", indentStr, fieldName, valueStr))
		}
	}

	return nil
}

// local helper for YAML generator
func isSimpleKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		return true
	default:
		return false
	}
}

// generateCommentFromTags генерирует комментарий из структурных тегов
func generateCommentFromTags(tag reflect.StructTag) string {
	var comments []string

	// Получаем info тег для основного описания
	info := tag.Get("info")
	if info != "" {
		comments = append(comments, info)
	}

	// Парсим RST теги
	tagsList := parseStructTag(tag)

	// Добавляем комментарии в детерминированном порядке
	ordered := []tagName{RST_MIN, RST_MAX, RST_DEFAULT, RST_CHOICE, RST_FORBIDDEN, RST_REGEX}
	for _, tn := range ordered {
		if tv, ok := tagsList[tn]; ok {
			comment := generateCommentForTag(tn, tv)
			if comment != "" {
				comments = append(comments, comment)
			}
		}
	}

	if len(comments) == 0 {
		return ""
	}

	return strings.Join(comments, "; ")
}

// generateCommentForTag генерирует комментарий для конкретного тега
func generateCommentForTag(tagName tagName, tagValue tagValue) string {
	switch tagName {
	case RST_MIN:
		return fmt.Sprintf("минимальное значение - %s", tagValue)

	case RST_MAX:
		return fmt.Sprintf("максимальное значение - %s", tagValue)

	case RST_DEFAULT:
		return fmt.Sprintf("значение по умолчанию - %s", tagValue)

	case RST_CHOICE:
		choices := strings.Split(string(tagValue), SET_DELIMITER)
		return fmt.Sprintf("допустимые значения: %s", strings.Join(choices, ", "))

	case RST_FORBIDDEN:
		parts := strings.Split(string(tagValue), VAL_DELIMITER)
		if len(parts) == 2 {
			// Есть список запрещенных значений + подменное значение
			forbidden := strings.Split(parts[0], SET_DELIMITER)
			repl := parts[1]
			return fmt.Sprintf("запрещенные значения: %s, подменное значение: %s",
				strings.Join(forbidden, ", "), repl)
		} else {
			// Только отдельные запрещенные значения
			forbidden := strings.Split(string(tagValue), SET_DELIMITER)
			return fmt.Sprintf("запрещенные значения: %s", strings.Join(forbidden, ", "))
		}

	case RST_REGEX:
		return fmt.Sprintf("регулярное выражение: %s", tagValue)

	default:
		return ""
	}
}

// formatValue форматирует значение для YAML
func formatValue(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf("\"%s\"", value.String())

	case reflect.Bool:
		return strconv.FormatBool(value.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)

	default:
		return fmt.Sprintf("%v", value.Interface())
	}
}
