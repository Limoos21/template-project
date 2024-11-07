package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed code/*
var templateFiles embed.FS

func readCodeFromEmbed(filePath, projectName string) (string, error) {
	data, err := templateFiles.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл %s: %w", filePath, err)
	}
	return replaceProjectName(string(data), projectName), nil
}

func replaceProjectName(code, projectName string) string {
	return strings.ReplaceAll(code, "{{ProjectName}}", projectName)
}

func initGoModule(projectName string) error {
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectName // Указываем директорию проекта
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("не удалось инициализировать go.mod: %w", err)
	}
	return nil
}

func tidyGoModule(projectName string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName // Указываем директорию проекта
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("не удалось выполнить go mod tidy: %w", err)
	}
	return nil
}

func main() {

	projectName := flag.String("name", "yourProjectName", "имя проекта")
	flag.Parse() // Парсим аргументы командной строки

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

	for filePath, codeSource := range files {
		content, err := readCodeFromEmbed(codeSource, *projectName)
		if err != nil {
			fmt.Printf("Не удалось прочитать встроенный шаблон %s для %s: %v\n", codeSource, filePath, err)
			return
		}

		if err := os.WriteFile(filePath, []byte(content), os.ModePerm); err != nil {
			fmt.Printf("Не удалось создать файл %s: %v\n", filePath, err)
			return
		}
	}

	if err := initGoModule(*projectName); err != nil {
		fmt.Printf("Ошибка инициализации go.mod: %v\n", err)
		return
	}

	if err := tidyGoModule(*projectName); err != nil {
		fmt.Printf("Ошибка выполнения go mod tidy: %v\n", err)
		return
	}

	fmt.Println("Проект успешно сгенерирован и инициализирован!")
}
