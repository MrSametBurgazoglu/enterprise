package tests

import (
	"testing"
	"time"

	"ariga.io/atlas/sql/schema"
	"github.com/MrSametBurgazoglu/enterprise/migrate"
	"github.com/MrSametBurgazoglu/enterprise/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTransformFieldToAtlasColumn_Defaults(t *testing.T) {
	tests := []struct {
		name     string
		field    models.FieldI
		expected string
		hasDef   bool
	}{
		{
			name:     "Bool Default True",
			field:    models.BoolField("IsActive").Default(true),
			expected: "true",
			hasDef:   true,
		},
		{
			name:     "Bool Default False",
			field:    models.BoolField("IsActive").Default(false),
			expected: "false",
			hasDef:   true,
		},
		{
			name:     "Int Default",
			field:    models.IntField("Count").Default(42),
			expected: "42",
			hasDef:   true,
		},
		{
			name:     "Uint Default",
			field:    models.UintField("Serial").Default(123),
			expected: "123",
			hasDef:   true,
		},
		{
			name:     "Float32 Default",
			field:    models.Float32Field("Price").Default(12.34),
			expected: "12.34",
			hasDef:   true,
		},
		{
			name:     "Float64 Default",
			field:    models.Float64Field("Price64").Default(56.78),
			expected: "56.78",
			hasDef:   true,
		},
		{
			name:     "String Default",
			field:    models.StringField("Status").Default("active"),
			expected: "'active'",
			hasDef:   true,
		},
		{
			name:     "String Default Escaping Single Quotes",
			field:    models.StringField("Name").Default("O'Connor"),
			expected: "'O''Connor'",
			hasDef:   true,
		},
		{
			name:     "Enum Default",
			field:    models.EnumField("DenemeType", []string{"Test", "Deneme"}).Default("Deneme"),
			expected: "'Deneme'",
			hasDef:   true,
		},
		{
			name:     "UUID DefaultFunc",
			field:    models.UUIDField("ID").DefaultFunc(uuid.New),
			expected: "gen_random_uuid()",
			hasDef:   true,
		},
		{
			name:     "Time DefaultFunc",
			field:    models.TimeField("CreatedAt").DefaultFunc(time.Now),
			expected: "now()",
			hasDef:   true,
		},
		{
			name:     "Byte Default",
			field:    models.ByteField("Data").Default([]byte{0xDE, 0xAD, 0xBE, 0xEF}),
			expected: "'\\xdeadbeef'",
			hasDef:   true,
		},
		{
			name:   "No Default",
			field:  models.StringField("Description"),
			hasDef: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			col := migrate.TransformFieldToAtlasColumn(tc.field)
			assert.NotNil(t, col)
			if tc.hasDef {
				assert.NotNil(t, col.Default)
				literal, ok := col.Default.(*schema.Literal)
				if assert.True(t, ok, "expected *schema.Literal") {
					assert.Equal(t, tc.expected, literal.V)
				}
			} else {
				assert.Nil(t, col.Default)
			}
		})
	}
}
