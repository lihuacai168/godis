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

var typeCmd = &cobra.Command{
	Use:   "type [key]",
	Short: "get key type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := typeKey(args[0])
		fmt.Printf("%s type is: %s\n", args[0], res)
	},
}
var ttlCmd = &cobra.Command{
	Use:   "ttl [key]",
	Short: "get key ttl",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := ttlKey(args[0])
		fmt.Println(res)
	},
}
var renamenxCmd = &cobra.Command{
	Use:     "renamenx [old_key] [new_key]",
	Aliases: []string{"mv"},
	Short:   "rename key, if new_key is exist return fail, else success",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//修改成功时，返回 1 。 如果 NEW_KEY_NAME 已经存在，返回 0 。
		var res bool
		if workedClient == "cluster" {
			res, _ = clusterClient.RenameNX(ctx, args[0], args[1]).Result()
		} else if workedClient == "alone" {
			res, _ = aloneClient.RenameNX(ctx, args[0], args[1]).Result()
		}

		if res {
			fmt.Println("rename success")
		} else {
			fmt.Println("rename fail")
		}
	},
}

func typeKey(key string) string {
	var r string
	if workedClient == "cluster" {
		r, _ = clusterClient.Type(ctx, key).Result()
	} else if workedClient == "alone" {
		r, _ = aloneClient.Type(ctx, key).Result()
	}
	return r
}

func ttlKey(key string) string {
	var r string
	if workedClient == "cluster" {
		r = clusterClient.TTL(ctx, key).String()
	} else if workedClient == "alone" {
		r = aloneClient.TTL(ctx, key).String()
	}
	return r
}
func isExists(key string) bool {
	var r int64
	if workedClient == "cluster" {
		r, _ = clusterClient.Exists(ctx, key).Result()
	} else if workedClient == "alone" {
		r, _ = aloneClient.Exists(ctx, key).Result()
	}
	return resultInt2Bool(r)
}

func resultInt2Bool(resCode int64) bool {
	if resCode == 1 {
		return true
	} else {
		return false
	}
}
func keys(pattern string) []string {
	result, _ := clusterClient.Keys(ctx, pattern).Result()
	return result
}
func init() {
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(existsCmd)
	rootCmd.AddCommand(typeCmd)
	rootCmd.AddCommand(renamenxCmd)
	rootCmd.AddCommand(ttlCmd)

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
