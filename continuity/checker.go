package continuity

import (
	"fmt"
	"maps"

	"github.com/verygoodsoftwarenotvirus/typewizard/models"
	"github.com/verygoodsoftwarenotvirus/typewizard/utils"
)

type mode int

const (
	Identical mode = iota
	LeftInclusive
)

type PackageDescription struct {
	Name string
	Path string
}

type Discrepancy struct {
	TypeName   string
	Reason     string
	PackageA   string
	PackageB   string
	FieldDiffs []FieldDifference
}

type FieldDifference struct {
	FieldName     string
	PackageAValue string
	PackageBValue string
	Property      string
}

func ComparePackageTypes(m mode, packages ...*PackageDescription) (bool, error) {
	if len(packages) <= 1 {
		return false, fmt.Errorf("must provide at least 2 packages")
	}

	typesMaps := make([]models.ListCollection[*models.Struct], len(packages))

	packageIndex := 0
	for _, pd := range packages {
		types, err := utils.GetTypesForPackage(pd.Path, pd.Name, nil)
		if err != nil {
			return false, err
		}

		typesMaps[packageIndex] = types.AsListCollection().Sort(func(x, y *models.Struct) bool {
			return x.Name < y.Name
		})

		packageIndex++
	}

	switch m {
	case LeftInclusive:
		if len(packages) != 2 {
			return false, fmt.Errorf("left inclusive mode only supports 2 packages")
		}

		var discrepancies []Discrepancy

		pkgA := typesMaps[0].AsMapCollection(func(m *models.Struct) string {
			return m.Name
		})
		pkgB := typesMaps[1].AsMapCollection(func(m *models.Struct) string {
			return m.Name
		})

		for typeName, structA := range pkgA {
			if structB, exists := pkgB[typeName]; !exists {
				discrepancies = append(discrepancies, Discrepancy{
					TypeName: typeName,
					Reason:   fmt.Sprintf("type missing in %s", typeName),
				})
				continue
			} else {
				if diff := compareStructs(structA, structB); diff != nil {
					discrepancies = append(discrepancies, *diff)
				}
			}
		}

		return len(discrepancies) == 0, nil
	case Identical:
		for i := range typesMaps {
			if i == 0 {
				continue
			}

			return typesMaps[i].EqualTo(typesMaps[0], func(structX, structY *models.Struct) bool {
				return structX.Fields.EqualTo(structY.Fields, func(fieldX, fieldY *models.StructField) bool {
					result := fieldX.Equal(fieldY)
					return result
				})
			}), nil
		}

		return false, nil
	default:
		return false, fmt.Errorf("unknown mode: %d", m)
	}
}

func compareStructs(a, b *models.Struct) *Discrepancy {
	var (
		diff       Discrepancy
		fieldDiffs []FieldDifference
	)
	aFields := a.Fields.AsMapCollection(func(f *models.StructField) string { return f.Name })
	bFields := b.Fields.AsMapCollection(func(f *models.StructField) string { return f.Name })

	for fieldName, fieldA := range aFields {
		fieldB, exists := bFields[fieldName]
		if !exists {
			fieldDiffs = append(fieldDiffs, FieldDifference{
				FieldName:     fieldName,
				Property:      "existence",
				PackageAValue: "exists",
				PackageBValue: "missing",
			})
			continue
		}

		if !fieldA.Equal(fieldB) {
			diffs := compareFields(fieldA, fieldB)
			fieldDiffs = append(fieldDiffs, diffs...)
		}
	}

	for fieldName := range bFields {
		if _, exists := aFields[fieldName]; !exists {
			fieldDiffs = append(fieldDiffs, FieldDifference{
				FieldName:     fieldName,
				Property:      "existence",
				PackageAValue: "missing",
				PackageBValue: "exists",
			})
		}
	}

	if len(fieldDiffs) > 0 {
		diff = Discrepancy{
			TypeName:   a.Name,
			Reason:     "field mismatch",
			FieldDiffs: fieldDiffs,
		}
		return &diff
	}

	return nil
}

func compareFields(a, b *models.StructField) []FieldDifference {
	var diffs []FieldDifference
	compareProperty := func(propName string, aVal, bVal interface{}) {
		aStr := fmt.Sprintf("%v", aVal)
		bStr := fmt.Sprintf("%v", bVal)
		if aStr != bStr {
			diffs = append(diffs, FieldDifference{
				FieldName:     a.Name,
				Property:      propName,
				PackageAValue: aStr,
				PackageBValue: bStr,
			})
		}
	}

	compareProperty("type", a.Type, b.Type)
	compareProperty("typePackage", a.TypePackage, b.TypePackage)
	compareProperty("typePackagePath", a.TypePackagePath, b.TypePackagePath)
	compareProperty("basicType", a.BasicType, b.BasicType)
	compareProperty("stdLib", a.FromStandardLibrary, b.FromStandardLibrary)

	if !maps.Equal(a.Tags, b.Tags) {
		compareProperty("tags", fmt.Sprintf("%v", a.Tags), fmt.Sprintf("%v", b.Tags))
	}

	return diffs
}
