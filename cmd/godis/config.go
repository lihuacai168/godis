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
	"github.com/lihuacai168/godis/cmd/config"
	"github.com/mattn/go-colorable"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	outWriter io.Writer = os.Stdout
	errWriter io.Writer = os.Stderr
	inReader  io.Reader = os.Stdin

	colorableOut io.Writer = colorable.NewColorableStdout()
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "Handle kaf configuration",
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAddClusterCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configRemoveClusterCmd)
	configCmd.AddCommand(configLsCmd)
	configCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.godis/config)")
	configAddClusterCmd.Flags().StringSliceVarP(&addrs, "addrs", "a", nil, "Comma separated list of addrs ip:port pairs")
	configAddClusterCmd.Flags().StringVarP(&password, "password", "p", "", "Redis password")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var configAddClusterCmd = &cobra.Command{
	Use:     "add-cluster [NAME]",
	Aliases: []string{"add"},
	Short:   "Add cluster",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		for _, cluster := range cfg.Clusters {
			if cluster.Name == name {
				errorExit("Could not add cluster: cluster with name '%v' exists already.", name)
			}
		}
		cfg.Clusters = append(cfg.Clusters, &config.Cluster{
			Name:     name,
			Addrs:    addrs,
			Password: password,
		})
		err := cfg.Write()
		if err != nil {
			errorExit("Unable to write config: %v\n", err)
		}
		fmt.Println("Added cluster.")
	},
}
var configRemoveClusterCmd = &cobra.Command{
	Use:               "remove-cluster [NAME]",
	Aliases:           []string{"rm", "delete", "del"},
	Short:             "remove cluster",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validConfigArgs,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		var pos = -1
		for i, cluster := range cfg.Clusters {
			if cluster.Name == name {
				pos = i
				break
			}
		}

		if pos == -1 {
			errorExit("Could not delete cluster: cluster with name '%v' not exists.", name)
		}

		cfg.Clusters = append(cfg.Clusters[:pos], cfg.Clusters[pos+1:]...)

		err := cfg.Write()
		if err != nil {
			errorExit("Unable to write config: %v\n", err)
		}
		fmt.Println("Removed cluster.")
	},
}
var configUseCmd = &cobra.Command{
	Use:               "use-cluster [NAME]",
	Aliases:           []string{"use"},
	Short:             "Sets the current cluster in the configuration",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validConfigArgs,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := cfg.SetCurrentCluster(name); err != nil {
			fmt.Printf("Cluster with name %v not found\n", name)
		} else {
			fmt.Printf("Switched to cluster \"%v\".\n", name)
		}
	},
}
var configLsCmd = &cobra.Command{
	Use:     "get-clusters",
	Aliases: []string{"ls", "list", "ll"},
	Short:   "Display clusters in the configuration file",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for _, cluster := range cfg.Clusters {
			marshal, _ := json.Marshal(cluster)
			os.Stdout.Write(marshal)
			fmt.Println()
		}
		fmt.Println("CurrentCluster: " + cfg.CurrentCluster)
	},
}

func validConfigArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	clusterList := make([]string, 0, len(cfg.Clusters))
	for _, cluster := range cfg.Clusters {
		clusterList = append(clusterList, cluster.Name)
	}
	return clusterList, cobra.ShellCompDirectiveNoFileComp
}
func errorExit(format string, a ...interface{}) {
	fmt.Fprintf(errWriter, format+"\n", a...)
	os.Exit(1)
}
