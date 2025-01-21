package helper

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
)

// CreateLocalEnv creates a .astra/env/env.yaml file
// Useful for commands that require this file and cannot create one on their own, for e.g. url, list
func CreateLocalEnv(context, compName, projectName string) {
	var config = fmt.Sprintf(`
ComponentSettings:
  Name: %s
  Project: %s
  AppName: app
`, compName, projectName)
	dir := filepath.Join(context, ".astra", "env")
	MakeDir(dir)
	Expect(os.WriteFile(filepath.Join(dir, "env.yaml"), []byte(config), 0600)).To(BeNil())
}
