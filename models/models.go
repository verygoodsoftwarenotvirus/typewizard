package models

import (
	"maps"
	"sort"
)

type Package struct {
	Path  string
	Alias string
	Name  string
}

type Struct struct {
	Name    string
	Package Package
	Fields  ListCollection[*StructField]
}

type StructField struct {
	Tags                map[string]string
	Name                string
	TypePackagePath     string
	TypePackage         string
	Type                string
	BasicType           bool
	FromStandardLibrary bool
}

func (f *StructField) Equal(f2 *StructField) bool {
	result := f.Name == f2.Name &&
		f.TypePackagePath == f2.TypePackagePath &&
		f.TypePackage == f2.TypePackage &&
		f.BasicType == f2.BasicType &&
		f.FromStandardLibrary == f2.FromStandardLibrary &&
		f.Type == f2.Type &&
		maps.Equal(f.Tags, f2.Tags)

	return result
}

type ListCollection[T any] []T

func (l *ListCollection[T]) AsMapCollection(keyFunc func(T) string) MapCollection[string, T] {
	out := make(MapCollection[string, T])
	for _, v := range *l {
		out[keyFunc(v)] = v
	}
	return out
}

func (l ListCollection[T]) Sort(lessFunc func(x, y T) bool) ListCollection[T] {
	sort.Slice(l, func(i, j int) bool {
		return lessFunc((l)[i], (l)[j])
	})

	return l
}

func (l *ListCollection[T]) EqualTo(l2 ListCollection[T], equalityFunc func(x, y T) bool) bool {
	if len(*l) != len(l2) {
		return false
	}

	for i, v := range *l {
		if i >= len(l2) {
			return false
		}
		if !equalityFunc(v, (l2)[i]) {
			return false
		}
	}
	return true
}

type MapCollection[K comparable, V any] map[K]V

func (m *MapCollection[K, V]) AsListCollection() ListCollection[V] {
	out := ListCollection[V]{}
	for _, v := range *m {
		out = append(out, v)
	}
	return out
}
