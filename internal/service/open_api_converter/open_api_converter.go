package open_api_converter

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"os/exec"
	"strings"
)

func ConvertOpenAPI(ctx context.Context, scriptDir, spec string) (map[string]any, error) {
	cmd := exec.CommandContext(ctx, "node", "convert.js")
	cmd.Dir = scriptDir
	cmd.Stdin = strings.NewReader(spec)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("openapi-to-postman: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	var collection map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &collection); err != nil {
		return nil, fmt.Errorf("decode collection: %w", err)
	}
	return collection, nil
}
