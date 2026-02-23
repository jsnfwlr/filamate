// Package database provides commands to manage the database.
package database

import "github.com/spf13/cobra"

var BaseCmd = &cobra.Command{
	Use:   "database",
	Short: "manage filamate's database",
}
