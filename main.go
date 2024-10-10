package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bbshebang",
	Short: "Manage your shebang line",
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

		cmd := exec.Command(command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting command: %v\n", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
		}
	} else {
		fmt.Println("Error reading shebang line.")
	}
}

func main() {
	rootCmd.AddCommand(runCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
