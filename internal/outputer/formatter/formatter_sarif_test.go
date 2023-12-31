package formatter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test"
	"github.com/qri-io/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestFormatSarif(t *testing.T) {
	sample := scheme_test.SchemeSample()

	for _, f := range []bool{true, false} {
		bytes, err := formatter.Format(formatter.Sarif, formatter.DefaultOutputIndent, sample, f)
		require.Nilf(t, err, "Error formatting sarif: %v", err)
		require.NotNil(t, bytes, "Error formatting sarif")
		require.NotEmpty(t, bytes, "Error formatting sarif")

		ctx := context.Background()

		schemaData, err := os.ReadFile("formatter_test/sarif_v2.1.0_schema.json")
		require.Nil(t, err)

		// QRI + JSON schema draft-07 compatibility
		// See https://github.com/qri-io/jsonschema/issues/114#issuecomment-1102010496
		jsonschema.RegisterKeyword("definitions", jsonschema.NewDefs)

		rs := &jsonschema.Schema{}
		err = json.Unmarshal(schemaData, rs)
		require.Nilf(t, err, "unmarshal schema: %v", err)

		errs, err := rs.ValidateBytes(ctx, bytes)
		require.Nil(t, err)

		if len(errs) > 0 {
			fmt.Println(errs[0].Error())
		}

		require.Emptyf(t, errs, "SARIF output does not match schema: %v", errs)
	}
}
