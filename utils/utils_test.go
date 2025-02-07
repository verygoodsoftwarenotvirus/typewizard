package utils

import (
	"testing"

	"github.com/verygoodsoftwarenotvirus/typewizard/models"

	"github.com/stretchr/testify/assert"
)

func Test_getTypesForPackage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := models.MapCollection[string, *models.Struct]{
			"ThingOne": &models.Struct{
				Name: "ThingOne",
				Fields: models.ListCollection[*models.StructField]{
					{
						Name:      "Name",
						Type:      "string",
						BasicType: true,
					},
					{
						Name:      "Age",
						Type:      "int",
						BasicType: true,
					},
					{
						Name:                "DateOfBirth",
						TypePackage:         "time",
						TypePackagePath:     "time",
						Type:                "Time",
						FromStandardLibrary: true,
					},
					{
						Name:            "Else",
						TypePackage:     "c",
						TypePackagePath: "typewizard/utils/test_packages/example/c",
						Type:            "SomethingElse",
					},
				},
			},
		}

		actual, err := GetTypesForPackage("./test_packages/example/a", "a", func(string) bool {
			return true
		})
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
