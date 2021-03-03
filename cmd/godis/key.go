/*
Copyright © 2021 lihuacai

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// keyCmd represents the key command
var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("key called")
	},
}

var deleteCmd = &cobra.Command{
	Use:     "del [key]",
	Aliases: []string{"rm", "delete"},
	Short:   "delete a key",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !isExists(args[0]) {
			fmt.Println("key not exists")
		} else {
			result := delete(args[0])
			var outMsg string
			if result == 1 {
				outMsg = "delete key success"
			} else {
				outMsg = "delete key fail"
			}
			_, _ = colorableOut.Write([]byte(outMsg))
			fmt.Fprintln(outWriter)
		}
	},
}

var existsCmd = &cobra.Command{
	Use:   "exists [key]",
	Short: "assure a key is exists",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if isExists(args[0]) {
			fmt.Println("key exists")
		} else {
			fmt.Println("key not exists")
		}
	},
}
var keysCmd = &cobra.Command{
	Use:   "keys [pattern]",
	Short: "search pattern keys",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := keys(args[0])
		for _, patternKey := range res {
			fmt.Println(patternKey)
		}
	},
}

func isExists(key string) bool {
	result, _ := rdb.Exists(ctx, key).Result()
	return resultInt2Bool(result)
}

func resultInt2Bool(resCode int64) bool {
	if resCode == 1 {
		return true
	} else {
		return false
	}
}
func keys(pattern string) []string {
	result, _ := rdb.Keys(ctx, pattern).Result()
	return result
}
func init() {
	rootCmd.AddCommand(keyCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(existsCmd)
	// 返回结果不准确，暂时注释
	//rootCmd.AddCommand(keysCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
