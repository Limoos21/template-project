package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// readCodeFromFile считывает содержимое файла и возвращает его как строку с заменой {{ProjectName}} на имя проекта
func readCodeFromFile(filePath, projectName string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	// Заменяем {{ProjectName}} на фактическое имя проекта
	return replaceProjectName(string(data), projectName), nil
}

// replaceProjectName заменяет {{ProjectName}} на фактическое имя проекта
func replaceProjectName(code, projectName string) string {
	return strings.ReplaceAll(code, "{{ProjectName}}", projectName)
}

// initGoModule инициализирует go.mod для проекта
func initGoModule(projectName string) error {
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectName // Указываем директорию проекта
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("не удалось инициализировать go.mod: %w", err)
	}
	return nil
}

// tidyGoModule выполняет go mod tidy
func tidyGoModule(projectName string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName // Указываем директорию проекта
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("не удалось выполнить go mod tidy: %w", err)
	}
	return nil
}

func main() {
	// Определяем флаг для имени проекта
	projectName := flag.String("name", "yourProjectName", "имя проекта")
	flag.Parse() // Парсим аргументы командной строки

	// Определяем структуру проекта и файлы с кодом
	files := map[string]string{
		filepath.Join(*projectName, "cmd", "main.go"):                                   "code/main.txt",
		filepath.Join(*projectName, "internal", "domain", "entity.go"):                  "code/entity.txt",
		filepath.Join(*projectName, "internal", "application", "example_usecase.go"):    "code/example_usecase.txt",
		filepath.Join(*projectName, "internal", "interfaces", "http", "handler.go"):     "code/handler.txt",
		filepath.Join(*projectName, "internal", "infra", "db", "db.go"):                 "code/db.txt",
		filepath.Join(*projectName, "internal", "infra", "db", "example_repository.go"): "code/example_repository.txt",
		filepath.Join(*projectName, "internal", "infra", "middleware", "cors.go"):       "code/cors.txt",
		filepath.Join(*projectName, "internal", "infra", "middleware", "auth.go"):       "code/auth.txt",
		filepath.Join(*projectName, "internal", "infra", "router.go"):                   "code/router.txt",
		filepath.Join(*projectName, "config", "config.go"):                              "code/config.txt",
		filepath.Join(*projectName, ".env"):                                             "code/.env.txt",
		filepath.Join(*projectName, "dockerfile"):                                       "code/dockerfile.txt",
		filepath.Join(*projectName, "makefile"):                                         "code/makefile.txt",
	}

	// Создаем директории
	dirs := []string{
		filepath.Join(*projectName, "cmd"),
		filepath.Join(*projectName, "internal", "domain"),
		filepath.Join(*projectName, "internal", "application"),
		filepath.Join(*projectName, "internal", "interfaces", "http"),
		filepath.Join(*projectName, "internal", "infra", "db"),
		filepath.Join(*projectName, "internal", "infra", "middleware"),
		filepath.Join(*projectName, "config"),
		filepath.Join(*projectName, "pkg"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Не удалось создать папку %s: %v\n", dir, err)
			return
		}
	}

	// Создаем файлы и добавляем содержимое из внешних файлов
	for filePath, codeSource := range files {
		var content string
		var err error

		// Если код для файла указан, считываем его из внешнего файла
		if codeSource != "" {
			content, err = readCodeFromFile(codeSource, *projectName) // Передаем имя проекта
			if err != nil {
				fmt.Printf("Не удалось прочитать файл %s для %s: %v\n", codeSource, filePath, err)
			}
		} else {
			content = "// Пустой файл, добавь свой код"
		}

		// Создаем файл и записываем код
		if err := os.WriteFile(filePath, []byte(content), os.ModePerm); err != nil {
			fmt.Printf("Не удалось создать файл %s: %v\n", filePath, err)
			return
		}
	}

	// Инициализируем go.mod
	if err := initGoModule(*projectName); err != nil {
		fmt.Printf("Ошибка инициализации go.mod: %v\n", err)
		return
	}

	// Выполняем go mod tidy
	if err := tidyGoModule(*projectName); err != nil {
		fmt.Printf("Ошибка выполнения go mod tidy: %v\n", err)
		return
	}

	fmt.Println("Проект успешно сгенерирован и инициализирован!")
}
