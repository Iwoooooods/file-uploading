package squirtle

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type QueryConfig struct {
	Table          string   `yaml:"table"`
	QueryFilePaths []string `yaml:"query_file"`
}

type QueryConfigStore []QueryConfig

type QueryMapper map[string]string

const (
	DefaultQueryStoreConfigPath = "./config/querystore.yaml"
)

func InitalizeQueryStore(cfgFilePath ...string) (QueryConfigStore, error) {
	if len(cfgFilePath) == 0 {
		cfgFilePath = append(cfgFilePath, DefaultQueryStoreConfigPath)
	}

	byt, err := os.ReadFile(cfgFilePath[0])
	if err != nil {
		return nil, err
	}

	var store QueryConfigStore
	if err := yaml.Unmarshal(byt, &store); err != nil {
		return nil, err
	}

	return store, nil
}

func (store QueryConfigStore) HydrateQueryStore(table string) (QueryMapper, error) {
	var config *QueryConfig
	var mapper = make(QueryMapper)

	for _, cfg := range store {
		if cfg.Table == table {
			config = &cfg
			break
		}
	}

	if config == nil {
		return nil, fmt.Errorf("no config found for table: %s", table)
	}

	if len(config.QueryFilePaths) == 0 {
		return nil, fmt.Errorf("no query file paths found for table: %s", table)
	}

	byt, err := os.ReadFile(config.QueryFilePaths[0])
	if err != nil {
		return mapper, err
	}

	queries := string(byt)

	rx := regexp.MustCompile(`(?m)sql:([P<QueryName>\w]+)$`)

	for _, queryConfig := range strings.Split(queries, "--") {
		matches := rx.FindStringSubmatch(queryConfig)
		if len(matches) < 2 {
			continue
		}

		mapper[matches[1]] = strings.TrimSpace(rx.ReplaceAllString(queryConfig, ""))
	}

	return mapper, nil
}

func (qm QueryMapper) Keys() []string {
	keys := make([]string, 0, len(qm))
	for k := range qm {
		keys = append(keys, k)
	}
	fmt.Println(keys)
	return keys
}

func (qm QueryMapper) GetQueries() []string {
	queries := make([]string, 0, len(qm))
	for _, v := range qm {
		queries = append(queries, v)
	}
	return queries
}

func (qm QueryMapper) GetQuery(queryName string) (string, error) {
	query, ok := qm[queryName]
	if !ok {
		return "", fmt.Errorf("query not found: %s", queryName)
	}
	fmt.Println("corresponding query:", query)
	return query, nil
}
