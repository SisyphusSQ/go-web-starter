//go:build integration

package scaf_fold

import "testing"

func TestGenerateE2EDBCombosIntegration(t *testing.T) {
	testGenerateE2EDBCombos(t, true)
}
