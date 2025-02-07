package utils

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"path/filepath"
	"slices"

	"github.com/verygoodsoftwarenotvirus/typewizard/models"

	"golang.org/x/tools/go/packages"
)

// GetTypesForPackage fetches type definitions with full import path resolution.
func GetTypesForPackage(packagePath, packageName string, nameFilter func(string) bool) (models.MapCollection[string, *models.Struct], error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedSyntax |
			packages.NeedFiles |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedImports,
		Dir:   filepath.Dir(packagePath),
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found")
	}

	output := models.MapCollection[string, *models.Struct]{}
	for _, pkg := range pkgs {
		if pkg.Name != packageName {
			continue
		}
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				switch t := n.(type) {
				case *ast.TypeSpec:
					if structType, ok := t.Type.(*ast.StructType); ok {
						if nameFilter != nil && !nameFilter(t.Name.Name) {
							return true
						}
						output[t.Name.Name] = &models.Struct{
							Name:   t.Name.Name,
							Fields: GetFieldsForStruct(pkg.TypesInfo, structType),
						}
					}
				}
				return true
			})
		}
	}

	return output, nil
}

// GetFieldsForStruct extracts fields with proper import path resolution.
func GetFieldsForStruct(typeInfo *types.Info, structType *ast.StructType) models.ListCollection[*models.StructField] {
	var structFields models.ListCollection[*models.StructField]

	for _, field := range structType.Fields.List {
		sf := &models.StructField{}
		if len(field.Names) > 0 {
			sf.Name = field.Names[0].Name
		}

		// Get type information from type-checker
		tv, ok := typeInfo.Types[field.Type]
		if !ok {
			continue
		}

		// Handle different type cases
		switch t := tv.Type.(type) {
		case *types.Named:
			sf.Type = t.Obj().Name()
			if pkg := t.Obj().Pkg(); pkg != nil {
				sf.TypePackage = pkg.Name()
				sf.TypePackagePath = pkg.Path()
				sf.FromStandardLibrary = isStdlibPkg(pkg.Path())
			}
		case *types.Basic:
			sf.Type = t.Name()
			sf.BasicType = true
		case *types.Pointer:
			if named, ok2 := t.Elem().(*types.Named); ok2 {
				sf.Type = "*" + named.Obj().Name()
				if pkg := named.Obj().Pkg(); pkg != nil {
					sf.TypePackage = pkg.Name()
					sf.TypePackagePath = pkg.Path()
					sf.FromStandardLibrary = isStdlibPkg(pkg.Path())
				}
			}
		}

		structFields = append(structFields, sf)
	}

	return structFields
}

var (
	stdLibPackages []string
)

func init() {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		log.Fatalf("reading standard library: %v", err)
	}

	for _, pkg := range pkgs {
		stdLibPackages = append(stdLibPackages, pkg.PkgPath)
	}
}

func isStdlibPkg(path string) bool {
	result := slices.Contains(stdLibPackages, path)
	return result
}
