package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "CLI App Example",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from CLI!")
		},
	}

	rootCmd.Execute()
}
