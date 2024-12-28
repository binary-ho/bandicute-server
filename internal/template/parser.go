package template

import (
	"encoding/json"
	"fmt"
	"os"
)

const Path = "internal/templates/"

func parseTemplateByFormat(template interface{}, templateFileName string) error {
	templateBytes, err := os.ReadFile(Path + templateFileName)
	if err != nil {
		return fmt.Errorf("failed to read summary prompt promptTemplate: %w", err)
	}

	if err := json.Unmarshal(templateBytes, &template); err != nil {
		return fmt.Errorf("failed to parseFeed summary prompt promptTemplate: %w", err)
	}
	return nil
}
