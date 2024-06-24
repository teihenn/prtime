package service

import (
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/teihenn/prtime/internal/model"
)

func TestGetSortedPRs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// モックのレスポンスを設定
	httpmock.RegisterResponder("GET", "https://api.bitbucket.org/2.0/repositories/testowner/testrepo/pullrequests",
		httpmock.NewStringResponder(200,
			`{"values":[{"id":1,"title":"Old PR","author":{"display_name":"user1"},
			"created_on":"2020-01-01T12:00:00Z"},
			{"id":2,"title":"New PR","author":{"display_name":"user2"},
			"created_on":"2022-01-01T12:00:00Z"}]}`))

	service := &BitbucketService{
		Repositories: []struct {
			Owner  string
			Name   string
			ApiKey string
		}{
			{Owner: "testowner", Name: "testrepo", ApiKey: "testapikey"},
		},
	}

	// DisplaySortedPRsを実行
	actual := service.GetSortedPRs()

	// ----- assert ----- //
	expected := map[string][]model.PullRequest{
		"testowner/testrepo": {
			{
				ID:        1,
				Title:     "Old PR",
				Author:    model.AuthorInfo{Username: "user1"},
				CreatedOn: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:        2,
				Title:     "New PR",
				Author:    model.AuthorInfo{Username: "user2"},
				CreatedOn: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}
