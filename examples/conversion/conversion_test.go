package conversion

import (
	"testing"

	"github.com/verygoodsoftwarenotvirus/typewizard/models"

	"github.com/stretchr/testify/assert"
)

func TestWriteConversionFunctionForTypes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		typeA := &models.Struct{
			Name: "ThingOne",
			Package: models.Package{
				Path:  "internal/services/handlers/models",
				Alias: "servicemodels",
				Name:  "servicemodels",
			},
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
			},
		}

		typeB := &models.Struct{
			Name: "ThingOne",
			Package: models.Package{
				Path:  "internal/database/models",
				Name:  "dbmodels",
				Alias: "dbmodels",
			},
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
			},
		}

		expected := `import (
	servicemodels "internal/services/handlers/models"
	dbmodels "internal/database/models"
)

func ConvertThingOneToThingOne(input servicemodels.ThingOne) dbmodels.ThingOne {
	out := dbmodels.ThingOne{}
	out.Name = input.Name
	out.Age = input.Age
	return out
}
`

		actual := WriteConversionFunctionForTypes(typeA, typeB)

		assert.Equal(t, expected, actual)
	})
}
