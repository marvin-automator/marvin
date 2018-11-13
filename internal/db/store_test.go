package db

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

type Point struct {
	X int
	Y int
}

func TestStore_Set_Get(t *testing.T) {
	SetupTestDB()
	defer TearDownTestDB()
	pt1 := Point{42, 88}
	var pt2 Point

	r := require.New(t)
	s := GetStore("myStore")
	s.Set("myPoint", pt1)
	s.Get("myPoint", &pt2)

	r.Equal(pt1, pt2)
}

func TestStore_Delete(t *testing.T) {
	SetupTestDB()
	defer TearDownTestDB()

	pt := Point{42, 42}

	r := require.New(t)
	s := GetStore("myStore")
	s.Set("WhatsThePoint", pt)
	s.Delete("WhatsThePoint")
	err := s.Get("WhatsThePoint", &pt)

	r.IsType(KeyNotFoundError{}, err)
}

func TestStore_EachKeyWithPrefix(t *testing.T) {
	SetupTestDB()
	defer TearDownTestDB()
	s1 := GetStore("store1")
	s2 := GetStore("store2")

	s2.Set("pt1", Point{1, 2})

	pts := []Point{
		{4, 3}, {8, 9}, {121, 45}, {90987, 123},
	}

	for i, pt := range pts {
		s1.Set("pt"+strconv.Itoa(i), pt)
	}

	s1.Set("notpt4", Point{404, 404})

	resultValues := make([]Point, 0, 4)
	keys := make([]string, 0, 4)

	var pt *Point = new(Point)
	err := s1.EachKeyWithPrefix("pt", pt, func(key string) error {
		keys = append(keys, key)
		resultValues = append(resultValues, *pt)
		return nil
	})

	r := require.New(t)
	r.NoError(err)
	r.Equal([]string{"pt0", "pt1", "pt2", "pt3"}, keys)
	r.Equal(pts, resultValues)
}
