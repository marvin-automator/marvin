package graphql

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"
)

func init() {
	RegisterTypeTransformer(func(t time.Time) string {
		return t.Format(time.RFC3339)
	})
}

type Profile struct {
	Created time.Time `json:"created"`
	Name string `json:"name"`
	LotteryNumbers []uint8 `json:"lottery_numbers"`
	Friends []Profile `json:"friends"`
	Buddy *Profile `json:"-"`
}

func isUnexporte(p cmp.Path) bool {
	fx, ok := p.Last().(cmp.StructField)
	if !ok {
		return false
	}
	r, _ := utf8.DecodeRuneInString(fx.Name())
	return unicode.IsLower(r)
}

func TestCreateOutputTypeFromStruct(t *testing.T) {
	r := require.New(t)
	ot := CreateOutputTypeFromStruct(Profile{})
	f := ot.(*graphql.Object).Fields()

	r.Equal(graphql.String, f["name"].Type)
	r.Equal(graphql.String, f["created"].Type)
	ignore := cmp.FilterPath(isUnexporte, cmp.Ignore())
	cmpopts.IgnoreTypes()
	r.Empty(cmp.Diff(graphql.NewList(graphql.Int), f["lottery_numbers"].Type, ignore))
	r.Empty(cmp.Diff(graphql.NewList(ot), f["friends"].Type, ignore))
	r.Equal(4, len(f))
}

func TestRegisterTypeTransformer(t *testing.T) {
	r := require.New(t)

	f := CreateOutputTypeFromStruct(Profile{}).(*graphql.Object).Fields()["created"]
	n := time.Now()
	res, err := f.Resolve(graphql.ResolveParams{
		Source: Profile{Created:n},
	})

	r.NoError(err)
	r.Equal(n.Format(time.RFC3339), res.(string))
}
