package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractTag(t *testing.T) {
	_, err := exractTag(nil)
	require.Error(t, err)

	meta := struct {
		ContentType string `rh:"Content-Type"`
		Token       string `rh:"Token,omitempty"`
		Page        int    `rq:"page"`
		CurrentPage int    `rq:"current_page, omitempty"`
	}{"json", "", 1, 0}

	_, err = exractTag(meta)
	require.NoError(t, err)

	extracted, err := exractTag(&meta)
	require.NoError(t, err)
	header := extracted[headerTagName]
	require.Equal(t, "json", header["Content-Type"])
	require.NotContains(t, header, "Token")

	query := extracted[queryTagName]
	require.Equal(t, "1", query["page"])
	require.NotContains(t, query, "current_page")

}
