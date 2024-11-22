package squirtle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitalizeQueryStore(t *testing.T) {
	store, err := InitalizeQueryStore("../../config/querystore.yaml")
	if err != nil {
		t.Fatal(err)
	}

	for _, cfg := range store {
		t.Log(cfg.Table)
		for _, qf := range cfg.QueryFilePaths {
			t.Log(qf)
		}
	}
}

func TestHydrateQueryStore(t *testing.T) {
	store := QueryConfigStore{
		{
			Table: "users",
			QueryFilePaths: []string{
				"./queries.sql",
			},
		},
	}
	_, err := store.HydrateQueryStore("lol")
	require.Error(t, err, "should error to get non-existent query mapper")

	qm, err := store.HydrateQueryStore("users")
	require.NoError(t, err, "should not error to get users query mapper")

	require.Equal(t, 2, len(qm.Keys()), "should have 2 queries")

	query, err := qm.GetQuery("CreateUserQuery")
	require.NoError(t, err, "should not error to get query")
	require.NotEmpty(t, query, "query should not be empty")
}
