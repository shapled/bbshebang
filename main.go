package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

const version = "0.1.1"

var rootCmd = &cobra.Command{
	Use:   "bbshebang",
	Short: "Manage your shebang line",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of bbshebang",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("bbshebang version %s\n", version)
	},
}

var runCmd = &cobra.Command{
	Use:   "run [filename]",
	Short: "Run a babashka script specified by the filename(need bb)",
	Args:  cobra.ExactArgs(1),
	Run:   runScript,
}

func runScript(cmd *cobra.Command, args []string) {
	filename := args[0]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		shebang := scanner.Text()
		if !strings.HasPrefix(shebang, "#!") {
			fmt.Println("No shebang found.")
			return
		}

		command := strings.TrimSpace(strings.TrimPrefix(shebang, "#!"))
		if runtime.GOOS == "windows" {
			if strings.HasPrefix(command, "/usr/bin/env ") {
				command = strings.TrimPrefix(command, "/usr/bin/env ")
			} else if strings.HasPrefix(command, "/usr/bin/") {
				command = strings.TrimPrefix(command, "/usr/bin/")
			} else if strings.HasPrefix(command, "/bin/") {
				command = strings.TrimPrefix(command, "/bin/")
			}
		}

		r := strings.NewReader(fmt.Sprintf("%s %s", command, filename))
		f, err := syntax.NewParser().Parse(r, "")
		if err != nil {
			fmt.Printf("Error parsing shebang: %v\n", err)
			return
		}

		syntax.NewPrinter().Print(os.Stdout, f)

		runner, err := interp.New(
			interp.StdIO(nil, os.Stdout, os.Stdout),
		)
		if err != nil {
			fmt.Printf("Error starting command: %v\n", err)
			return
		}

		if err := runner.Run(context.TODO(), f); err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
		}
	} else {
		fmt.Println("Error reading shebang line.")
	}
}

func main() {
	rootCmd.AddCommand(runCmd, versionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
