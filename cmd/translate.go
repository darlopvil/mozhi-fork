package cmd

import (
	"encoding/json"
	"strings"
	"fmt"

	"codeberg.org/aryak/libmozhi"

	"github.com/spf13/cobra"
)

var (
	engine  string
	engines  string
	query   string
	source  string
	dest    string
	rawjson bool
)

func printEngineResult(result libmozhi.LangOut, printPlaceHolderAndEngineName bool) {
	if printPlaceHolderAndEngineName {
		fmt.Println("-----------------------------------")
		fmt.Println("Engine: " + result.Engine)
	}

	fmt.Println("Translated Text: " + result.OutputText)
	if source == "auto" {
		fmt.Println("Detected Language: " + result.AutoDetect)
	}
	fmt.Println("Source Language: " + result.SourceLang)
	fmt.Println("Target Language: " + result.TargetLang)
}

func printRaw(data interface{}) {
	j, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(j))
	}
}

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate.",
	Run: func(cmd *cobra.Command, args []string) {
		if engine == "all" {
			data := libmozhi.TranslateAll(dest, source, query)
			if rawjson {
				printRaw(data)
			} else {
				for _, result := range data {
					printEngineResult(result, true)
				}
				fmt.Println("-----------------------------------")
			}
		} else if engine == "some" {
			if engines != "" {
				data, err := libmozhi.TranslateSome(strings.Split(engines, ","), dest, source, query)
				if rawjson {
					printRaw(data)
				} else {
					if err != nil {
						fmt.Println(err)
					} else {
						for _, result := range data {
							printEngineResult(result, true)
						}
						fmt.Println("-----------------------------------")
					}
				}
			} else {
				fmt.Println("Please add the --engines / -E argument when using the engine `some`")
			}
		} else {
			data, err := libmozhi.Translate(engine, dest, source, query)
			if rawjson {
				printRaw(data)
			} else {
				if err != nil {
					fmt.Println(err)
				} else {
					printEngineResult(data, false)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().SortFlags = false

	translateCmd.Flags().StringVarP(&engine, "engine", "e", "", "[all|some|google|libre|reverso|deepl|yandex|mymemory|duckduckgo]")
	translateCmd.Flags().StringVarP(&engines, "engines", "E", "", "Engines to select. This flag is to be used with the `some` engine only")
	translateCmd.Flags().StringVarP(&source, "source", "s", "", "Source language. Use langlist command to get code for your language")
	translateCmd.Flags().StringVarP(&dest, "dest", "t", "", "Target language. Use langlist command to get code for your language")
	translateCmd.Flags().StringVarP(&query, "query", "q", "", "Text to be translated")
	translateCmd.Flags().BoolVarP(&rawjson, "raw", "r", false, "Return output as json")

	translateCmd.MarkFlagRequired("source")
	translateCmd.MarkFlagRequired("dest")
	translateCmd.MarkFlagRequired("query")
}
