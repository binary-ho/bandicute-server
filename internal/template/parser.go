package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const templatesDirectoryName = "templates"

func parseTemplateByFormat(template interface{}, templateFileName string) error {
	templatePath, err := getTemplateFilePath(templateFileName)
	if err != nil {
		return fmt.Errorf("failed to get template file path: %w", err)
	}

	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read summary prompt promptTemplate: %w", err)
	}

	if err := json.Unmarshal(templateBytes, template); err != nil {
		return fmt.Errorf("failed to parseFeed summary prompt promptTemplate: %w", err)
	}
	return nil
}

func getTemplateFilePath(templateFileName string) (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get runtime caller")
	}

	basePath := filepath.Dir(filename)
	return filepath.Join(basePath, templatesDirectoryName, templateFileName), nil
}
