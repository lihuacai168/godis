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
	"encoding/json"
	"fmt"
	"github.com/hokaccha/go-prettyjson"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set key operation",
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(SMembersCmdShort)
	setCmd.AddCommand(SMembersCmd)

	rootCmd.AddCommand(SAddCmdShort)
	setCmd.AddCommand(SAddCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var SMembersCmdShort = SMembersCmd
var SMembersCmd = &cobra.Command{
	Use:     "smembers [key]",
	Aliases: []string{"sget"},
	Short:   "get set key values",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result := SMembers(args[0])
		_, _ = colorableOut.Write(result)
		fmt.Fprintln(outWriter)
	},
}

func SMembers(key string) []byte {
	result, _ := rdb.SMembers(ctx, key).Result()
	sl := make([]interface{}, len(result))
	for i, s := range result {
		m, err := JsonToMap(s)
		if err != nil {
			sl[i] = s
		} else {
			sl[i] = m
		}
	}
	b, _ := json.Marshal(sl)
	p, _ := prettyjson.Format(b)
	return p
}

var SAddCmdShort = SAddCmd
var SAddCmd = &cobra.Command{
	Use:   "sadd [key] member1 member2 member3...",
	Short: "add a set key",
	Args:  cobra.RangeArgs(2, int(^uint(0)>>1)),
	Run: func(cmd *cobra.Command, args []string) {
		members := make([]interface{}, len(args))
		for index, arg := range args {
			if index == 0 {
				continue
			}
			members[index] = arg
		}
		SAdd(args[0], members)
		result := SMembers(args[0])
		fmt.Println("sadd success, set key is "+args[0], ", value is ")
		_, _ = colorableOut.Write(result)
		fmt.Fprintln(outWriter)
	},
}

func SAdd(key string, members []interface{}) {
	rdb.SAdd(ctx, key, members...)
}
