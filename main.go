package main

import (
	"fmt"

	adapter "adapter/adapt"
)

func main() {
	fmt.Println("Демонстрация пакета Adapt")
	fmt.Println("=========================")
	fmt.Println()

	// Демонстрация основного функционала AdaptStruct
	fmt.Println("=== ДЕМОНСТРАЦИЯ ADAPTSTRUCT ===")
	adapter.RunAdaptExamples()
	fmt.Println()

	// Демонстрация YAML генератора
	fmt.Println("=== ДЕМОНСТРАЦИЯ YAML ГЕНЕРАТОРА ===")

	// Простой пример
	fmt.Println("1. Простой пример:")
	fmt.Println("------------------")
	adapter.RunExample()
	fmt.Println()

	// Сложный пример
	fmt.Println("2. Сложный пример с вложенными структурами:")
	fmt.Println("--------------------------------------------")
	adapter.RunComplexExample()
	fmt.Println()

	// Пример с ошибкой
	fmt.Println("3. Пример с ошибкой (не структура):")
	fmt.Println("-----------------------------------")
	_, err := adapter.GenerateStructYAML("не структура")
	if err != nil {
		fmt.Printf("Ожидаемая ошибка: %v\n", err)
	}
	fmt.Println()

	// Пример создания файла конфигурации
	fmt.Println("4. Создание файла конфигурации:")
	fmt.Println("--------------------------------")
	adapter.RunFileOnlyExample()
	fmt.Println()

	fmt.Println("Демонстрация завершена!")
}
