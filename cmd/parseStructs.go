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
	"encoding/json"
	"flag"
	"fmt"
	"github.com/RivenZoo/goext/code_parser"
)

var parseStructCfg parseStructConfig

type parseStructConfig struct {
	sourceFilePath string
}

var subCmd = &subCommand{
	Name:  "parse_struct",
	flags: flag.NewFlagSet("parse_struct_fs", flag.ContinueOnError),
	Usage: "Parse struct info from go source code.",
	Execute: func() error {
		fileInfo, err := code_parser.ParseSourceCode(parseStructCfg.sourceFilePath)
		if err == nil {
			b, _ := json.MarshalIndent(fileInfo, "", "  ")
			fmt.Println(string(b))
		}
		return err
	},
}

func init() {
	subCmd.flags.StringVar(&parseStructCfg.sourceFilePath, "srcFile", "", "source code file path")
	rootCmd.addSubCommand(subCmd)
}
