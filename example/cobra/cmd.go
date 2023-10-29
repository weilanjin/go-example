package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// go get -u github.com/spf13/cobra

// cobra-cli 模版生成
// go install github.com/spf13/cobra-cli@latest

var rootCmd = &cobra.Command{
	Use:   "go",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains`,
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("import called")
		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("export called")
		return nil
	},
}

func main() {
	rootCmd.AddCommand(importCmd, exportCmd)
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}
