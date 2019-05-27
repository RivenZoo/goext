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
package code_parser

/*
  code_parser parse go source code and extract code token information.
*/

type CommentGroup []string

type ImportInfo struct {
	// import alias
	Alias      string `json:"alias" yaml:"alias"`
	ImportPath string `json:"import_path" yaml:"import_path"` // import module path
}

type StructFieldInfo struct {
	Name             string            `json:"name" yaml:"name"`
	Type             string            `json:"type" yaml:"type"`
	IsEmbbedType     bool              `json:"is_embbed_type" yaml:"is_embbed_type"`
	IsAnonymousField bool              `json:"is_anonymous_field" yaml:"is_anonymous_field"`
	AnonymousStruct  *StructDefinition `json:"anonymous_struct" yaml:"anonymous_struct"`
	Docs             CommentGroup      `json:"docs" yaml:"docs"`
	Comments         CommentGroup      `json:"comments" yaml:"comments"`
	Tag              string            `json:"tag" yaml:"tag"`
}

type StructDefinition struct {
	Fields []StructFieldInfo `json:"fields" yaml:"fields"`
}

type TypeDefineInfo struct {
	Name        string            `json:"name" yaml:"name"`
	IsTypeAlias bool              `json:"is_type_alias" yaml:"is_type_alias"`
	SourceType  string            `json:"source_type" yaml:"source_type"`
	Docs        CommentGroup      `json:"docs" yaml:"docs"`
	Definition  *StructDefinition `json:"definition" yaml:"definition"`
}

// SourceFileInfo contain file level information like package name, import, struct definition and comments.
type SourceFileInfo struct {
	SourceFilePath  string           `json:"source_file_path" yaml:"source_file_path"`
	PackageName     string           `json:"package_name" yaml:"package_name"`
	ImportPackages  []ImportInfo     `json:"import_packages" yaml:"import_packages"`
	Docs            CommentGroup     `json:"docs" yaml:"docs"`
	TypeDefinitions []TypeDefineInfo `json:"type_definitions" yaml:"type_definitions"`
}
