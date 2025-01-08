package cmd

import (
	"fmt"
	"log"
	"os"

	"android-translation-tool/cmd/internal"
	"android-translation-tool/cmd/ui/singleselect"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("key", "k", "", "Key to use for translation (no spaces allowed, lowercases letters and underscores only)")
	createCmd.Flags().StringP("value", "v", "", "String to translate (english only, closed in quotes)")
	createCmd.Flags().BoolP("apply", "a", false, "Apply the translation to the project (default is false) (if false it will only print the translations)")
}

var createCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate a string",
	Run: func(cmd *cobra.Command, args []string) {
		inAndroidProject := internal.CheckCurrentDirectoryIsAndroidProject()
		if !inAndroidProject {
			fmt.Println("This is not an Android project or you are not in the root directory of an Android project.")
			return
		}

		key := cmd.Flag("key").Value.String()
		str := cmd.Flag("value").Value.String()

		fmt.Println("Key:", key)
		fmt.Println("String:", str)
		fmt.Println("")

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}
		resDirs := internal.FindResourcesDirectoriesPath(currentDir)

		if len(resDirs) == 0 {
			fmt.Println("No Android resource directories found.")
			return
		}

		selectedPath := singleselect.Selection{Selected: ""}

		var tprogram *tea.Program
		tprogram = tea.NewProgram(singleselect.InitialModelSingleSelect(resDirs, &selectedPath))
		if _, err := tprogram.Run(); err != nil {
			log.Printf("Name of project contains an error: %v", err)
		}

		if selectedPath.Selected == "" {
			return
		}

		strings := internal.GetTranslationsFromResourceDirectory(selectedPath.Selected)

		languagesFound := []string{}
		for _, s := range strings {
			languagesFound = append(languagesFound, s.Language)
		}

		fmt.Println("Languages found:", languagesFound)

		fmt.Println("Translating...\n")

		for _, s := range strings {
			t, err := internal.TranslateText(str, s.LocaleCode)
			if err != nil {
				fmt.Println("Error translating text:", err)
				fmt.Println("Language:", s)
				return
			}

			fmt.Println(s.Language, ":", t)
		}
	},
}
