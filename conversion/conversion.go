package conversion

import (
	"fmt"
	"strings"

	"github.com/verygoodsoftwarenotvirus/typewizard/models"
)

func WriteConversionFunctionForTypes(typeA, typeB *models.Struct) string {
	var sb strings.Builder

	imports := []models.Package{
		typeA.Package,
		typeB.Package,
	}

	var importBlock strings.Builder
	if len(imports) > 0 {
		importBlock.WriteString("import (\n")
		for _, pkg := range imports {
			if pkg.Alias != "" {
				importBlock.WriteString(fmt.Sprintf("\t%s \"%s\"\n", pkg.Alias, pkg.Path))
			} else {
				importBlock.WriteString(fmt.Sprintf("\t\"%s\"\n", pkg.Path))
			}
		}
		importBlock.WriteString(")\n\n")
	}

	sb.WriteString(importBlock.String())

	funcName := fmt.Sprintf("Convert%sTo%s", typeA.Name, typeB.Name)
	sb.WriteString(fmt.Sprintf("func %s(input %s.%s) %s.%s {\n", funcName, typeA.Package.Name, typeA.Name, typeB.Package.Name, typeB.Name))
	sb.WriteString(fmt.Sprintf("\tout := %s.%s{}\n", typeB.Package.Name, typeB.Name))

	for _, fieldA := range typeA.Fields {
		for _, fieldB := range typeB.Fields {
			if fieldA.Name == fieldB.Name {
				if fieldA.Type == fieldB.Type {
					sb.WriteString(fmt.Sprintf("\tout.%s = input.%s\n", fieldB.Name, fieldA.Name))
				} else if (fieldA.Type == "float32" && fieldB.Type == "float64") || (fieldA.Type == "float64" && fieldB.Type == "float32") {
					sb.WriteString(fmt.Sprintf("\tout.%s = %s(input.%s)\n", fieldB.Name, fieldB.Type, fieldA.Name))
				}
			}
		}
	}

	sb.WriteString("\treturn out\n")
	sb.WriteString("}\n")

	return sb.String()
}
