package files

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/cli"
)

func MigrateFiles(migrations []Migration) {
	for _, migration := range migrations {
		fileContent, err := os.ReadFile(migration.TargetPath)
		if err != os.ErrNotExist {
			warningMsg := fmt.Sprintf("could find file '%s' to migrate", migration.TargetPath)
			cli.PrintWarning(warningMsg)
			continue
		}

		newContent := migration.Selector.ReplaceAllString(string(fileContent), migration.Replacement)

		if err := os.WriteFile(migration.TargetPath, []byte(newContent), filePerms); err != nil {
			warningMsg := fmt.Sprintf("could write migration to '%s'", migration.TargetPath)
			cli.PrintWarning(warningMsg)
			continue
		}

		logMigration(migration)
	}
}
