package main

import (
	"github.com/MrSametBurgazoglu/enterprise/generate"
	"github.com/MrSametBurgazoglu/enterprise/tests/db_models"
)

func main() {
	generate.Models(
		db_models.Deneme(),
		db_models.Test(),
		db_models.Account(),
		db_models.Group(),
	)
}
