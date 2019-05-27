// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

var rootCmd = rootCommand{
	subCommands: map[string]*subCommand{},
}

type subCommand struct {
	flags   *flag.FlagSet
	Usage   string
	Name    string
	Execute func() error
}

type rootCommand struct {
	subCommands map[string]*subCommand
}

func (r *rootCommand) addSubCommand(c *subCommand) {
	if _, ok := r.subCommands[c.Name]; ok {
		panic(fmt.Errorf("sub command %s exists", c.Name))
	}
	r.subCommands[c.Name] = c
}

func (r *rootCommand) usage() string {
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	buf.WriteString("Usage: goext command [option]\n\nCommands:\n")
	for name, c := range r.subCommands {
		buf.WriteString(fmt.Sprintf("  %s:\n\t", name))
		buf.WriteString(c.Usage)
		buf.WriteString("\n")
		w := c.flags.Output()
		c.flags.SetOutput(buf)
		c.flags.PrintDefaults()
		c.flags.SetOutput(w)
	}
	buf.WriteString("\n")
	return buf.String()
}

func (r *rootCommand) parseFlag() (cmdName string) {
	if len(os.Args) < 2 {
		fmt.Println(r.usage())
		os.Exit(-1)
	}
	cmdName = os.Args[1]

	var c *subCommand
	var ok bool
	if c, ok = r.subCommands[cmdName]; !ok {
		fmt.Println(r.usage())
		os.Exit(-2)
	}

	if c.flags != nil {
		if len(os.Args) < 3 {
			fmt.Println(r.usage())
			os.Exit(-1)
		}
		err := c.flags.Parse(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse flag error %v", err)
			os.Exit(-1)
		}
	}
	return
}

func (r *rootCommand) execute() {
	cmdName := rootCmd.parseFlag()
	c := r.subCommands[cmdName]
	if err := c.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "execute %s error %v", cmdName, err)
		fmt.Println(r.usage())
		os.Exit(-2)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.execute()
}

func init() {
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}
