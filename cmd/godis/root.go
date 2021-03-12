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
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lihuacai168/godis/cmd/config"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	outWriter io.Writer = os.Stdout
	errWriter io.Writer = os.Stderr
	inReader  io.Reader = os.Stdin

	colorableOut   io.Writer = colorable.NewColorableStdout()
	ConnectSuccess           = "PONG"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godis",
	Short: "A utility redis command line",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		outWriter = cmd.OutOrStdout()
		errWriter = cmd.ErrOrStderr()
		inReader = cmd.InOrStdin()
		if outWriter != os.Stdout {
			colorableOut = outWriter
		}
		if currentCluster != nil && currentCluster.IsSafeMode {
			parentCmd := cmd.Parent().Use
			if parentCmd != "config" {
				nameAndAliases := strings.Split(cmd.NameAndAliases(), ",")
				safeCmd := []string{"hg", "hget", "hgetall", "hga", "sget", "smembers", "type", "config", "ttl"}
				if !contains(safeCmd, nameAndAliases[0]) {
					panic(fmt.Sprintf("safe mode can only support cmds: %v", safeCmd))
				}
			}
		}

	},
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var cfg config.Config
var workedClient string
var currentCluster *config.Cluster
var clusterClient *redis.ClusterClient
var aloneClient *redis.Client
var ctx = context.Background()
var (
	password           string
	clusterDescription string
	addrs              []string
	clusterOverride    string
	isSafeMode         bool
)

func init() {
	cobra.OnInitialize(onInit)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&clusterOverride, "cluster", "c", "", "set a temporary current cluster")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".godis" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".godis")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
func onInit() {
	var err error
	cfg, err = config.ReadConfig(cfgFile)
	if err != nil {
		errorExit("Invalid config: %v", err)
	}

	cfg.ClusterOverride = clusterOverride

	cluster := cfg.ActiveCluster()
	if cluster != nil {
		// Use active cluster from config
		currentCluster = cluster
	}
	initClient()
}
func initClient() {
	if currentCluster == nil {
		if cfg.ClusterOverride != "" {
			panic(fmt.Sprintf("%s not in configuration", cfg.ClusterOverride))
		} else {
			log.Println(`not configured, please use "godis config add" to set a cluster configuration`)
		}
	} else {
		a := currentCluster.Addrs
		p := currentCluster.Password
		clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    a,
			Password: p,
		})
		clusterSuccess, err := clusterClient.Ping(ctx).Result()
		aloneClient = redis.NewClient(&redis.Options{Addr: a[0], Password: p})
		aloneSuccess, _ := aloneClient.Ping(ctx).Result()
		if clusterSuccess == ConnectSuccess {
			workedClient = "cluster"
		} else if aloneSuccess == ConnectSuccess {
			workedClient = "alone"
		}

		if workedClient == "" {
			log.Printf("cluster and alone mode connect failed, use %s conf, addr: %s, password: %s\nPing server: %v\n", currentCluster.Name, currentCluster.Addrs, currentCluster.Password, err)
		} else {
			log.Printf("connect success, using %s mode, conf: %s, addr is %s \n", workedClient, currentCluster.Name, a)
		}
	}

}
