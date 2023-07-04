package main

import (
	"testing"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/codecoverage"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/stretchr/testify/require"
)

func TestRetrieveCoverageReports(t *testing.T) {
	testSession := session.NewSession()

	client, err := ranger.NewClient("", testSession)
	require.NoError(t, err)

	err = codecoverage.KillAgentTestServicesRetrieveCoverage(client)
	require.NoError(t, err)

	err = codecoverage.KillRangerTestServicesRetrieveCoverage(client)
	require.NoError(t, err)

}
