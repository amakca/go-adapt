package adapter

import (
	"fmt"
)

// ExampleStruct демонстрирует использование различных RST тегов
type ExampleStruct struct {
	// Простые поля с различными тегами
	Counter     int     `json:"counter" rst-min:"5" info:"Счетчик"`
	Price       float64 `rst-max:"1000.0" info:"Цена товара"`
	Status      string  `rst-choice:"active||inactive||pending" info:"Статус заказа"`
	Email       string  `rst-regex:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$" info:"Email адрес"`
	UserID      uint    `rst-forbidden:"0||1||2**10" info:"ID пользователя"`
	Description string  `rst-default:"Без описания" info:"Описание"`

	// Сложные поля
	Tags     []string               `rst-max:"10" info:"Теги"`
	Settings map[string]interface{} `info:"Настройки"`
	IsActive *bool                  `rst-default:"true" info:"Активен ли"`
}

// RunExample демонстрирует работу функции GenerateStructYAML
func RunExample() {
	// Создаем экземпляр структуры
	example := ExampleStruct{
		Counter:     10,
		Price:       99.99,
		Status:      "active",
		Email:       "user@example.com",
		UserID:      5,
		Description: "Пример товара",
		Tags:        []string{"tag1", "tag2", "tag3"},
		Settings: map[string]interface{}{
			"theme": "dark",
			"lang":  "ru",
		},
		IsActive: nil,
	}

	// Генерируем YAML и выводим в консоль
	yaml, err := GenerateStructYAML(example)
	if err != nil {
		fmt.Printf("Ошибка при генерации YAML: %v\n", err)
		return
	}

	fmt.Println("Сгенерированный YAML:")
	fmt.Println("=====================")
	fmt.Println(yaml)

	// Создаем YAML файл
	err = GenerateStructYAMLFile(example, "example_struct")
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		return
	}

	fmt.Println("\nФайл example_struct.yaml создан успешно!")
}

// Пример использования с вложенной структурой
type User struct {
	Name   string `json:"name" rst-regex:"[a-zA-Z]+" info:"Имя пользователя"`
	Age    int    `rst-min:"18" rst-max:"120" info:"Возраст"`
	Email  string `rst-regex:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$" info:"Email"`
	Active bool   `rst-default:"true" info:"Активен ли"`
}

type Order struct {
	ID       int      `rst-min:"1" info:"Номер заказа"`
	User     User     `info:"Информация о пользователе"`
	Items    []string `rst-max:"20" info:"Список товаров"`
	Total    float64  `rst-min:"0.0" rst-max:"10000.0" info:"Общая сумма"`
	Status   string   `rst-choice:"new||processing||completed||cancelled" info:"Статус заказа"`
	Discount float64  `rst-forbidden:"0.0||-1.0**-100.0" info:"Скидка"`
}

func RunComplexExample() {
	order := Order{
		ID: 12345,
		User: User{
			Name:   "Иван",
			Age:    25,
			Email:  "ivan@example.com",
			Active: true,
		},
		Items:    []string{"Товар 1", "Товар 2", "Товар 3"},
		Total:    1500.50,
		Status:   "processing",
		Discount: 10.0,
	}

	// Генерируем YAML и выводим в консоль
	yaml, err := GenerateStructYAML(order)
	if err != nil {
		fmt.Printf("Ошибка при генерации YAML: %v\n", err)
		return
	}

	fmt.Println("Сложный пример - Заказ:")
	fmt.Println("========================")
	fmt.Println(yaml)

	// Создаем YAML файл
	err = GenerateStructYAMLFile(order, "complex_order")
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		return
	}

	fmt.Println("\nФайл complex_order.yaml создан успешно!")
}

// RunFileOnlyExample демонстрирует создание YAML файла без вывода в консоль
func RunFileOnlyExample() {
	// Создаем структуру для конфигурации
	type Config struct {
		Server struct {
			Host string `json:"host" rst-default:"localhost" info:"Хост сервера"`
			Port int    `json:"port" rst-min:"1" rst-max:"65535" rst-default:"8080" info:"Порт сервера"`
		} `json:"server" info:"Настройки сервера"`

		Database struct {
			Driver   string `json:"driver" rst-choice:"mysql||postgres||sqlite" info:"Тип базы данных"`
			Host     string `json:"host" rst-default:"localhost" info:"Хост БД"`
			Port     int    `json:"port" rst-min:"1" rst-max:"65535" rst-default:"5432" info:"Порт БД"`
			Username string `json:"username" rst-regex:"[a-zA-Z0-9_]+" info:"Имя пользователя"`
			Password string `json:"password" info:"Пароль"`
		} `json:"database" info:"Настройки базы данных"`

		Logging struct {
			Level      string `json:"level" rst-choice:"debug||info||warn||error" rst-default:"info" info:"Уровень логирования"`
			File       string `json:"file" rst-default:"app.log" info:"Файл логов"`
			MaxSize    int    `json:"max_size" rst-min:"1" rst-max:"100" rst-default:"10" info:"Максимальный размер файла (МБ)"`
			MaxBackups int    `json:"max_backups" rst-min:"0" rst-max:"10" rst-default:"5" info:"Количество резервных копий"`
		} `json:"logging" info:"Настройки логирования"`
	}

	// Создаем экземпляр конфигурации
	config := Config{}
	config.Server.Host = "0.0.0.0"
	config.Server.Port = 9000
	config.Database.Driver = "postgres"
	config.Database.Host = "db.example.com"
	config.Database.Port = 5432
	config.Database.Username = "app_user"
	config.Database.Password = "secret_password"
	config.Logging.Level = "debug"
	config.Logging.File = "application.log"
	config.Logging.MaxSize = 20
	config.Logging.MaxBackups = 3

	// Создаем YAML файл
	err := GenerateStructYAMLFile(config, "config")
	if err != nil {
		fmt.Printf("Ошибка при создании файла конфигурации: %v\n", err)
		return
	}

	fmt.Println("Файл config.yaml создан успешно!")
	fmt.Println("Этот файл содержит конфигурацию с комментариями из структурных тегов.")
}
