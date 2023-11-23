/*
Copyright © 2023 bri <b@ibeep.com>

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
	"errors"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const AppName = "gomox"

var cfgFile string

var c context.Context

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gomox",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	var (
		pveUser     string
		pvePassword string
		pveRealm    string
		pveUrl      string
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli2.yaml)")

	rootCmd.PersistentFlags().StringVar(&pveUser, "pveuser", "root", "The username to log in as")
	viper.BindPFlag("pve_user", rootCmd.PersistentFlags().Lookup("pveuser"))
	viper.RegisterAlias("pveuser", "pve_user")

	rootCmd.PersistentFlags().StringVar(&pvePassword, "pvepassword", "root", "The password to log in with")
	viper.BindPFlag("pve_password", rootCmd.PersistentFlags().Lookup("pvepassword"))
	viper.RegisterAlias("pvepassword", "pve_password")
	rootCmd.PersistentFlags().StringVar(&pveRealm, "pverealm", "pam", "The realm to log in to")
	viper.BindPFlag("pverealm", rootCmd.PersistentFlags().Lookup("pverealm"))
	viper.RegisterAlias("pverealm", "pve_realm")
	rootCmd.PersistentFlags().StringVar(&pveUrl, "pveurl", "", "PVE URL")
	viper.BindPFlag("pveurl", rootCmd.PersistentFlags().Lookup("pveurl"))
	viper.RegisterAlias("pveurl", "pve_url")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	if len(cfgFile) > 0 { // Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if xdg := os.Getenv("XDG_CONFIG_HOME"); len(xdg) > 0 {
		viper.AddConfigPath(path.Join(xdg, AppName))
		viper.SetConfigName("config")
	} else {
		viper.AddConfigPath(path.Join(home, ".config", AppName))
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match
	err = viper.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if errors.As(err, &configFileNotFoundError) {
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("." + AppName)
		err = viper.ReadInConfig()
		if errors.As(err, &configFileNotFoundError) {
			viper.SetConfigType("envfile")
			viper.SetConfigName(".env")
			err = viper.ReadInConfig()
			if errors.As(err, &configFileNotFoundError) {
				err = nil // It's ok if the config file does not exist
			}
		}
	}
	cobra.CheckErr(err)
}
