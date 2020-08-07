package kustomize

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Build(dir string) (string, error) {
	var b strings.Builder
	c := exec.Command("kustomize", "build", dir)
	c.Stdout = &b
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return "", fmt.Errorf("could not execute kustomize: %w", err)
	}
	return b.String(), nil
}
