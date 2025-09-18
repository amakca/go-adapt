package adapter

import (
	"fmt"
	"log"
	"os"
)

// RunAdaptExamples демонстрирует работу основного функционала AdaptStruct
func RunAdaptExamples() {
	fmt.Println("Демонстрация работы AdaptStruct:")
	fmt.Println("================================")

	// Пример 1: Простая адаптация
	runSimpleAdaptExample()

	// Пример 2: Адаптация с выбором значений
	runChoiceAdaptExample()

	// Пример 3: Адаптация с запрещенными значениями
	runForbiddenAdaptExample()

	// Пример 4: Адаптация слайсов
	runSliceAdaptExample()

	// Пример 5: Адаптация вложенных структур
	runNestedAdaptExample()
}

// runSimpleAdaptExample демонстрирует простую адаптацию с min/max/default
func runSimpleAdaptExample() {
	fmt.Println("\n1. Простая адаптация (min/max/default):")

	type SimpleStruct struct {
		Count  int     `rst-min:"5" rst-max:"100" rst-default:"10"`
		Price  float64 `rst-min:"0.0" rst-max:"1000.0" rst-default:"99.99"`
		Name   string  `rst-default:"Unknown"`
		Active bool    `rst-default:"true"`
	}

	original := SimpleStruct{
		Count:  3,      // Меньше минимума
		Price:  1500.0, // Больше максимума
		Name:   "",     // Пустая строка
		Active: false,  // Будет изменено на true
	}

	fmt.Printf("До адаптации: %+v\n", original)

	var a = adapter{logger: log.New(os.Stdout, "adapter ", log.LstdFlags)}

	adapted, err := a.AdaptStruct(original)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("После адаптации: %+v\n", adapted)
}

// runChoiceAdaptExample демонстрирует адаптацию с выбором значений
func runChoiceAdaptExample() {
	fmt.Println("\n2. Адаптация с выбором значений (choice):")

	type ChoiceStruct struct {
		Status string `rst-choice:"active||inactive||pending"`
		Type   int    `rst-choice:"1||2||3"`
		Color  string `rst-choice:"red||green||blue"`
	}

	original := ChoiceStruct{
		Status: "unknown", // Неизвестный статус
		Type:   5,         // Неизвестный тип
		Color:  "yellow",  // Неизвестный цвет
	}

	fmt.Printf("До адаптации: %+v\n", original)

	var a = adapter{logger: log.New(os.Stdout, "adapter ", log.LstdFlags)}

	adapted, err := a.AdaptStruct(original)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("После адаптации: %+v\n", adapted)
}

// runForbiddenAdaptExample демонстрирует адаптацию с запрещенными значениями
func runForbiddenAdaptExample() {
	fmt.Println("\n3. Адаптация с запрещенными значениями (forbidden):")

	type ForbiddenStruct struct {
		ID    int     `rst-forbidden:"0||1||2**10"`
		Value float64 `rst-forbidden:"-1.0||0.0**-10.0"`
		Code  string  `rst-forbidden:"error||fail||invalid"`
	}

	original := ForbiddenStruct{
		ID:    5,       // В запрещенном диапазоне
		Value: -5.0,    // В запрещенном диапазоне
		Code:  "error", // Запрещенное значение
	}

	fmt.Printf("До адаптации: %+v\n", original)

	var a = adapter{logger: log.New(os.Stdout, "adapter ", log.LstdFlags)}

	adapted, err := a.AdaptStruct(original)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("После адаптации: %+v\n", adapted)
}

// runSliceAdaptExample демонстрирует адаптацию слайсов
func runSliceAdaptExample() {
	fmt.Println("\n4. Адаптация слайсов (max):")

	type SliceStruct struct {
		Numbers []int     `rst-max:"5"`
		Names   []string  `rst-max:"3"`
		Prices  []float64 `rst-max:"10"`
	}

	original := SliceStruct{
		Numbers: []int{1, 2, 3, 4, 5, 6, 7, 8},                // 8 элементов
		Names:   []string{"Alice", "Bob", "Charlie", "David"}, // 4 элемента
		Prices:  []float64{10.0, 20.0},                        // 2 элемента (в пределах лимита)
	}

	fmt.Printf("До адаптации:\n")
	fmt.Printf("  Numbers: %v (длина: %d)\n", original.Numbers, len(original.Numbers))
	fmt.Printf("  Names: %v (длина: %d)\n", original.Names, len(original.Names))
	fmt.Printf("  Prices: %v (длина: %d)\n", original.Prices, len(original.Prices))

	var a = adapter{logger: log.New(os.Stdout, "adapter ", log.LstdFlags)}

	adapted, err := a.AdaptStruct(original)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	adaptedStruct := adapted.(SliceStruct)
	fmt.Printf("После адаптации:\n")
	fmt.Printf("  Numbers: %v (длина: %d)\n", adaptedStruct.Numbers, len(adaptedStruct.Numbers))
	fmt.Printf("  Names: %v (длина: %d)\n", adaptedStruct.Names, len(adaptedStruct.Names))
	fmt.Printf("  Prices: %v (длина: %d)\n", adaptedStruct.Prices, len(adaptedStruct.Prices))
}

// runNestedAdaptExample демонстрирует адаптацию вложенных структур
func runNestedAdaptExample() {
	fmt.Println("\n5. Адаптация вложенных структур:")

	type Address struct {
		Street  string `rst-default:"Unknown Street"`
		City    string `rst-default:"Unknown City"`
		ZipCode string `rst-regex:"[0-9]{5}"`
	}

	type User struct {
		Name    string   `rst-regex:"[a-zA-Z]+" rst-default:"Unknown"`
		Age     int      `rst-min:"18" rst-max:"120" rst-default:"25"`
		Email   string   `rst-regex:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
		Address Address  `info:"Адрес пользователя"`
		Tags    []string `rst-max:"5"`
	}

	original := User{
		Name:  "John123",       // Содержит цифры
		Age:   15,              // Меньше минимума
		Email: "invalid-email", // Неверный формат
		Address: Address{
			Street:  "",    // Пустая строка
			City:    "",    // Пустая строка
			ZipCode: "123", // Неверный формат
		},
		Tags: []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6", "tag7"}, // 7 тегов
	}

	fmt.Printf("До адаптации:\n")
	fmt.Printf("  Name: %s\n", original.Name)
	fmt.Printf("  Age: %d\n", original.Age)
	fmt.Printf("  Email: %s\n", original.Email)
	fmt.Printf("  Address: %+v\n", original.Address)
	fmt.Printf("  Tags: %v (длина: %d)\n", original.Tags, len(original.Tags))

	var a = adapter{logger: log.New(os.Stdout, "adapter ", log.LstdFlags)}

	adapted, err := a.AdaptStruct(original)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	adaptedUser := adapted.(User)
	fmt.Printf("После адаптации:\n")
	fmt.Printf("  Name: %s\n", adaptedUser.Name)
	fmt.Printf("  Age: %d\n", adaptedUser.Age)
	fmt.Printf("  Email: %s\n", adaptedUser.Email)
	fmt.Printf("  Address: %+v\n", adaptedUser.Address)
	fmt.Printf("  Tags: %v (длина: %d)\n", adaptedUser.Tags, len(adaptedUser.Tags))
}
