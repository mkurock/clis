package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func main() {
  cmd := &cobra.Command{
    Use: "findandreplace [flags] directory",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      run(cmd.Flag("config").Value.String(), args[0])
    },
  }
  cmd.Flags().StringP("config", "c", "replace.yaml", "yaml file with replace contents")
  cmd.Execute()
}

func run(replaceYaml string, path string) {
	content, err := os.ReadFile(replaceYaml)
	if err != nil {
		log.Fatal("Could not open config file")
		return
	}
	var data map[string]map[string]string
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		log.Fatal("Error parsing yaml")
    return
	}
	for file, replaceMap := range data {
		fmt.Printf("parsing filepattern %v\n", file)
		fileResult, err := doublestar.Glob(os.DirFS(path), file)
		if err != nil {
			log.Fatal("could not process file glob")
      continue
		}
		for _, res := range fileResult {
      fmt.Printf("Found file: %v\n", res)
      fileContentBytes, err := os.ReadFile(filepath.Join(path, res))
      if err != nil {
        log.Fatalf("error reading file: %v\n", res)
        continue
      }
      fileContent := string(fileContentBytes)
			for before, after := range replaceMap {
        fileContent = strings.ReplaceAll(fileContent, before, after)
			}
      fmt.Println(fileContent)
      os.WriteFile(filepath.Join(path, res), []byte(fileContent), os.ModePerm)
		}
	}
}
