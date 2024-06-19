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
	expectedApiKeyPath := "expected/api/key/path"
	if actual.Bitbucket.ApiKeyPath != expectedApiKeyPath {
		t.Errorf("expected:%s, actual:%s", expectedApiKeyPath, actual.Bitbucket.ApiKeyPath)
	}
	expectedRepositories := []struct {
		Owner string
		Name  string
	}{
		{Owner: "expectedOwner1", Name: "expectedRepoName1"},
		{Owner: "expectedOwner2", Name: "expectedRepoName2"},
	}

	if len(actual.Bitbucket.Repositories) != len(expectedRepositories) {
		t.Errorf("expected %d repositories, got %d", len(expectedRepositories), len(actual.Bitbucket.Repositories))
	}
	for i, repo := range expectedRepositories {
		if actual.Bitbucket.Repositories[i].Owner != repo.Owner || actual.Bitbucket.Repositories[i].Name != repo.Name {
			t.Errorf("expected repository %d: owner=%s, name=%s; got owner=%s, name=%s",
				i, repo.Owner, repo.Name, actual.Bitbucket.Repositories[i].Owner, actual.Bitbucket.Repositories[i].Name)
		}
	}

}
