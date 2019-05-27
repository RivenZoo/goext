package code_parser

import (
	"github.com/RivenZoo/backbone/logger"
	"github.com/patrickmn/go-cache"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"time"
)

func ParseSourceCode(srcFile string) (*SourceFileInfo, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, srcFile, nil, parser.ParseComments)
	if err != nil {
		logger.Errorf("parse source %s error %v", srcFile, err)
		return nil, err
	}

	absSrcFile, err := filepath.Abs(srcFile)
	if err != nil {
		logger.Errorf("get abs path %s error %v", srcFile, err)
		return nil, err
	}
	absSrcFile = filepath.Clean(absSrcFile)

	importInfo := parseImports(fs, f)
	docs := parseFileDocs(fs, f)
	defineInfos := parseStructs(fs, f)

	srcInfo := &SourceFileInfo{
		SourceFilePath:  absSrcFile,
		PackageName:     f.Name.Name,
		ImportPackages:  importInfo,
		Docs:            docs,
		TypeDefinitions: defineInfos,
	}
	return srcInfo, nil
}

func parseImports(fs *token.FileSet, f *ast.File) []ImportInfo {
	ret := make([]ImportInfo, 0, len(f.Imports))
	for i := range f.Imports {
		imp := f.Imports[i]
		var alias, importPath string
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		if imp.Path != nil {
			importPath = imp.Path.Value
		}
		ret = append(ret, ImportInfo{
			Alias:      alias,
			ImportPath: importPath,
		})
	}
	return ret
}

func parseFileDocs(fs *token.FileSet, f *ast.File) CommentGroup {
	comments := make(CommentGroup, 0, 16)
	if f.Doc != nil {
		comments = append(comments, f.Doc.Text())
	}
	return comments
}

func parseStructs(fs *token.FileSet, f *ast.File) []TypeDefineInfo {
	structs := make([]TypeDefineInfo, 0, len(f.Decls))
	for i := range f.Decls {
		ast.Inspect(f.Decls[i], func(node ast.Node) bool {
			genDecl, ok := node.(*ast.GenDecl)
			if ok {
				if genDecl.Tok == token.TYPE {
					structInfo := parseStructDecl(fs, genDecl.Specs[0])
					if structInfo != nil {
						tpDoc := genDecl.Doc.Text()
						if tpDoc != "" {
							structInfo.Docs = append(structInfo.Docs, tpDoc)
						}

						structs = append(structs, *structInfo)
					}
				}
				return true
			}
			return true
		})
	}
	return structs
}

func parseStructDecl(fs *token.FileSet, node ast.Node) *TypeDefineInfo {
	tpSpec, ok := node.(*ast.TypeSpec)
	if !ok {
		return nil
	}
	if _, ok := tpSpec.Type.(*ast.InterfaceType); ok {
		// skip interface define
		return nil
	}
	structInfo := &TypeDefineInfo{
		Name:       tpSpec.Name.Name,
		Definition: &StructDefinition{},
	}

	switch structTp := tpSpec.Type.(type) {
	case *ast.StructType:
		structInfo.Definition.Fields = parseStructFields(fs, structTp.Fields)
	case *ast.InterfaceType:
	default:
		structInfo.IsTypeAlias = true
		structInfo.SourceType = readNodeToken(fs, structTp)
	}
	return structInfo
}

var sourceFileCache = cache.New(5*time.Minute, 10*time.Minute)

func readNodeToken(fs *token.FileSet, node ast.Node) string {
	return readTokenByPosition(fs.Position(node.Pos()), fs.Position(node.End()))
}

func readTokenByPosition(pos, end token.Position) string {
	data, found := sourceFileCache.Get(pos.Filename)

	var c []byte
	var err error
	if !found {
		c, err = ioutil.ReadFile(pos.Filename)
		if err != nil {
			logger.Errorf("read file %s error %v", pos.Filename, err)
			return ""
		}
		sourceFileCache.Set(pos.Filename, c, 0)
	} else {
		c = data.([]byte)
	}

	return string(c[pos.Offset:end.Offset])
}

func parseStructFields(fs *token.FileSet, node ast.Node) []StructFieldInfo {
	fields, ok := node.(*ast.FieldList)
	if !ok {
		return nil
	}
	fieldInfo := make([]StructFieldInfo, 0, len(fields.List))
	for i := range fields.List {
		field := fields.List[i]
		info := StructFieldInfo{}

		if field.Tag != nil {
			info.Tag = field.Tag.Value
		}

		if field.Doc != nil {
			info.Docs = CommentGroup{field.Doc.Text()}
		}

		if field.Comment != nil {
			info.Comments = CommentGroup{field.Comment.Text()}
		}

		if len(field.Names) == 0 {
			info.IsEmbbedType = true
		} else {
			info.Name = field.Names[0].Name
		}

		switch fieldTp := field.Type.(type) {
		case *ast.StructType:
			info.IsAnonymousField = true
			info.AnonymousStruct = &StructDefinition{
				Fields: parseStructFields(fs, fieldTp.Fields),
			}
		default:
			info.Type = readNodeToken(fs, field.Type)
		}
		fieldInfo = append(fieldInfo, info)
	}
	return fieldInfo
}
