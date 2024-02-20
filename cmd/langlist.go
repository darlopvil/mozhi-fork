package cmd

import (
	"fmt"

	"codeberg.org/aryak/libmozhi"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

var (
	engineused string
	listtype   string
	raw        bool
)

// langlistCmd represents the langlist command
var langlistCmd = &cobra.Command{
	Use:   "langlist",
	Short: "List/select languages supported by an engine.",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := libmozhi.LangList(engineused, listtype)
		if err != nil {
			fmt.Println(err)
		} else {
			idxs, err := fuzzyfinder.FindMulti(
				list,
				func(i int) string {
					return list[i].Name
				})
			if err != nil {
				fmt.Println(err)
			}
			for _, idx := range idxs {
				if raw {
					fmt.Println(list[idx].Id)
				} else {
					fmt.Println("Selected Language:", list[idx].Id, "("+list[idx].Name+")")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(langlistCmd)
	translateCmd.Flags().SortFlags = false

	langlistCmd.Flags().StringVarP(&engineused, "engine", "e", "", "[google|libre|reverso|deepl|watson|yandex|mymemory|duckduckgo]")
	langlistCmd.Flags().StringVarP(&listtype, "type", "t", "", "[sl|tl] Choose language for source or target")
	langlistCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Return only selected language code.")

	langlistCmd.MarkFlagRequired("engine")
	langlistCmd.MarkFlagRequired("type")
}
