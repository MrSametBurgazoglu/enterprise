package generate

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/MrSametBurgazoglu/enterprise/models"
	"os"
	"slices"
	"strings"
	"text/template"
)

//go:embed schema_struct.go.tpl
var schemaTemplate string

//go:embed predicates.go.tpl
var predicateTemplate string

//go:embed client.go.tpl
var clientTemplate string

func Models(tables ...*models.Table) {
	g := &models.Generation{}
	g.Tables = append(g.Tables, tables...)

	if err := os.Mkdir("models", 0755); err != nil && !os.IsExist(err) {
		panic(err)
	}

	Clients(g)
	Schemas(g)
	Predicates(g)
}

func Schemas(g *models.Generation) {
	t := template.Must(template.New("").Parse(schemaTemplate))
	for _, table := range g.Tables {
		for _, field := range table.Fields {
			for _, s := range field.GetRequiredPackages() {
				if !slices.Contains(table.RequiredPackages, s) {
					table.RequiredPackages = append(table.RequiredPackages, s)
				}
			}
		}
		buf := &bytes.Buffer{}
		filePath := fmt.Sprintf("models/%s.go", strings.ToLower(table.TableName))

		if err := t.Execute(buf, table); err != nil {
			panic(err)
		}
		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(buf.Bytes())
		if err != nil {
			panic(err)
		}
	}
}

func Predicates(g *models.Generation) {
	t := template.Must(template.New("").Parse(predicateTemplate))
	for _, table := range g.Tables {
		buf := &bytes.Buffer{}

		if err := t.Execute(buf, table); err != nil {
			panic(err)
		}

		filePath := fmt.Sprintf("models/%s_predicates.go", strings.ToLower(table.TableName))
		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(buf.Bytes())
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func Clients(g *models.Generation) {
	t := template.Must(template.New("").Parse(clientTemplate))

	buf := &bytes.Buffer{}

	if err := t.Execute(buf, g); err != nil {
		panic(err)
	}

	filePath := "models/client.go"
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
}
