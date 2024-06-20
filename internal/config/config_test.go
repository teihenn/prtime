package config

import (
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	/* ----- arrange ----- */
	testConfigPath := filepath.Join("testdata", "test_config.yml")

	/* ----- act ----- */
	actual, err := Load(testConfigPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	/* ----- assert ----- */
	expectedRepositories := []struct {
		Owner      string
		Name       string
		ApiKeyPath string
	}{
		{Owner: "expectedOwner1", Name: "expectedRepoName1", ApiKeyPath: "expected/api/key/path"},
		{Owner: "expectedOwner2", Name: "expectedRepoName2", ApiKeyPath: "expected/api/key/path"},
	}

	if len(actual.Bitbucket.Repositories) != len(expectedRepositories) {
		t.Errorf("expected %d repositories, got %d",
			len(expectedRepositories), len(actual.Bitbucket.Repositories))
	}
	for i, repo := range expectedRepositories {
		if actual.Bitbucket.Repositories[i].Owner != repo.Owner ||
			actual.Bitbucket.Repositories[i].Name != repo.Name {
			t.Errorf(
				"expected repository %d: owner=%s, name=%s, apiKeyPath=%s; "+
					"got owner=%s, name=%s, apiKeyPath=%s",
				i, repo.Owner, repo.Name, repo.ApiKeyPath,
				actual.Bitbucket.Repositories[i].Owner,
				actual.Bitbucket.Repositories[i].Name,
				actual.Bitbucket.Repositories[i].ApiKeyPath,
			)
		}
	}

}
