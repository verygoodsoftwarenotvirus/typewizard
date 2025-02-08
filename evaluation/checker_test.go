package evaluation

import (
	"testing"

	"github.com/verygoodsoftwarenotvirus/typewizard/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackagesHaveSameTypes_Identical(T *testing.T) {
	T.Parallel()

	T.Run("detects match", func(t *testing.T) {
		t.Parallel()

		input := []*PackageDescription{
			{
				Name: "b",
				Path: "test_packages/identical/good/b",
			},
			{
				Name: "b",
				Path: "test_packages/identical/good/b_copy",
			},
		}

		actual, err := ComparePackageTypes(
			Identical,
			input...,
		)
		assert.NoError(t, err)
		require.NotNil(t, actual)
		assert.True(t, actual)
	})

	T.Run("detects mismatch", func(t *testing.T) {
		t.Parallel()

		input := []*PackageDescription{
			{
				Name: "q",
				Path: "test_packages/identical/good/a",
			},
			{
				Name: "b",
				Path: "test_packages/identical/good/b",
			},
		}

		actual, err := ComparePackageTypes(
			Identical,
			input...,
		)
		assert.NoError(t, err)
		require.NotNil(t, actual)
		assert.False(t, actual)
	})
}

func TestPackagesHaveSameTypes_LeftInclusive(T *testing.T) {
	T.Parallel()

	T.Run("detects match", func(t *testing.T) {
		t.Parallel()

		input := []*PackageDescription{
			{
				Name: "a",
				Path: "test_packages/left_inclusive/good/a",
			},
			{
				Name: "b",
				Path: "test_packages/left_inclusive/good/b",
			},
		}

		actual, err := ComparePackageTypes(
			LeftInclusive,
			input...,
		)
		assert.NoError(t, err)
		require.NotNil(t, actual)
		assert.True(t, actual)
	})

	T.Run("detects mismatch", func(t *testing.T) {
		t.Parallel()

		input := []*PackageDescription{
			{
				Name: "a",
				Path: "test_packages/left_inclusive/bad/a",
			},
			{
				Name: "b",
				Path: "test_packages/left_inclusive/bad/b",
			},
		}

		actual, err := ComparePackageTypes(
			LeftInclusive,
			input...,
		)
		assert.NoError(t, err)
		require.NotNil(t, actual)
		assert.False(t, actual)
	})
}

func TestFindMatchingTypePairs(T *testing.T) {
	T.Parallel()

	T.Run("detects match", func(t *testing.T) {
		t.Parallel()

		pkgA, err := utils.GetTypesForPackage("test_packages/identical/good/b", "b", nil)
		require.NoError(t, err)
		require.NotNil(t, pkgA)

		pkgB, err := utils.GetTypesForPackage("test_packages/identical/good/b_copy", "b", nil)
		require.NoError(t, err)
		require.NotNil(t, pkgA)

		actual := FindMatchingTypePairs(pkgA, pkgB)
		assert.NoError(t, err)
		require.NotNil(t, actual)

		assert.Len(t, actual, 1)
		assert.Equal(t, actual[0].A.Name, pkgA.AsListCollection()[0].Name)
	})
}
