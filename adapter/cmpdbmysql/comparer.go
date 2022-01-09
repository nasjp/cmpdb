package cmpdbmysql

import (
	"github.com/nasjp/cmpdb"
)

func Load(cfg *cmpdb.Config) (*cmpdb.Comparer, error) {
	mysql := &mysql{db: cfg.DB}
	comparer := &cmpdb.Comparer{
		Adapter: mysql,
		Bytes:   cfg.Bytes,
	}

	if cfg.Bytes != nil {
		dbDiff, err := loadDiffJSON(cfg.Bytes)
		if err != nil {
			return nil, err
		}
		comparer.DBDiff = dbDiff

		if err := mysql.Load(dbDiff.BeforeDB); err != nil {
			return nil, err
		}
	}

	return comparer, nil
}

func loadDiffJSON(bytes []byte) (*cmpdb.DBDiff, error) {
	return cmpdb.ParseFromJSONDiff(bytes)
}
