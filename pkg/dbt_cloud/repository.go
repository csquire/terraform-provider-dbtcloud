package dbt_cloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Repository struct {
	ID                                    *int      `json:"id,omitempty"`
	AccountID                             int       `json:"account_id"`
	ProjectID                             int       `json:"project_id"`
	RemoteUrl                             string    `json:"remote_url"`
	State                                 int       `json:"state"`
	AzureActiveDirectoryProjectID         *string   `json:"azure_active_directory_project_id,omitempty"`
	AzureActiveDirectoryRepositoryID      *string   `json:"azure_active_directory_repository_id,omitempty"`
	AzureBypassWebhookRegistrationFailure *bool     `json:"azure_bypass_webhook_registration_failure,omitempty"`
	GitCloneStrategy                      string    `json:"git_clone_strategy"`
	RepositoryCredentialsID               *int      `json:"repository_credentials_id,omitempty"`
	GitlabProjectID                       *int      `json:"gitlab_project_id"`
	GithubInstallationID                  *int      `json:"github_installation_id"`
	DeployKey                             DeployKey `json:"deploy_key,omitempty"`
}

type DeployKey struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	State     int    `json:"state"`
	PublicKey string `json:"public_key"`
}

type RepositoryListResponse struct {
	Data   []Repository   `json:"data"`
	Status ResponseStatus `json:"status"`
}

type RepositoryResponse struct {
	Data   Repository     `json:"data"`
	Status ResponseStatus `json:"status"`
}

func (c *Client) GetRepository(
	repositoryID, projectID string,
) (*Repository, error) {

	repositoryUrl := fmt.Sprintf(
		"%s/v3/accounts/%s/projects/%s/repositories/%s/",
		c.HostURL,
		strconv.Itoa(c.AccountID),
		projectID,
		repositoryID,
	)

	req, err := http.NewRequest("GET", repositoryUrl, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	repositoryResponse := RepositoryResponse{}
	err = json.Unmarshal(body, &repositoryResponse)
	if err != nil {
		return nil, err
	}

	return &repositoryResponse.Data, nil
}

func (c *Client) CreateRepository(
	projectID int,
	remoteUrl string,
	isActive bool,
	gitCloneStrategy string,
	gitlabProjectID int,
	githubInstallationID int,
	azureActiveDirectoryProjectID string,
	azureActiveDirectoryRepositoryID string,
	azureBypassWebhookRegistrationFailure bool,
) (*Repository, error) {
	state := STATE_ACTIVE
	if !isActive {
		state = STATE_DELETED
	}

	newRepository := Repository{
		AccountID:        c.AccountID,
		ProjectID:        projectID,
		RemoteUrl:        remoteUrl,
		State:            state,
		GitCloneStrategy: gitCloneStrategy,
	}
	if gitlabProjectID != 0 {
		newRepository.GitlabProjectID = &gitlabProjectID
	}
	if githubInstallationID != 0 {
		newRepository.GithubInstallationID = &githubInstallationID
	}
	if azureActiveDirectoryProjectID != "" {
		newRepository.AzureActiveDirectoryProjectID = &azureActiveDirectoryProjectID
		newRepository.AzureActiveDirectoryRepositoryID = &azureActiveDirectoryRepositoryID
		newRepository.AzureBypassWebhookRegistrationFailure = &azureBypassWebhookRegistrationFailure
	}
	newRepositoryData, err := json.Marshal(newRepository)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/v3/accounts/%s/projects/%s/repositories/",
			c.HostURL,
			strconv.Itoa(c.AccountID),
			strconv.Itoa(projectID),
		),
		strings.NewReader(string(newRepositoryData)),
	)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	repositoryResponse := RepositoryResponse{}
	err = json.Unmarshal(body, &repositoryResponse)
	if err != nil {
		return nil, err
	}

	return &repositoryResponse.Data, nil
}

func (c *Client) UpdateRepository(
	repositoryID, projectID string,
	repository Repository,
) (*Repository, error) {
	repositoryData, err := json.Marshal(repository)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/v3/accounts/%s/projects/%s/repositories/%s/",
			c.HostURL,
			strconv.Itoa(c.AccountID),
			projectID,
			repositoryID,
		),
		strings.NewReader(string(repositoryData)),
	)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	repositoryResponse := RepositoryResponse{}
	err = json.Unmarshal(body, &repositoryResponse)
	if err != nil {
		return nil, err
	}

	return &repositoryResponse.Data, nil
}

func (c *Client) DeleteRepository(repositoryID, projectID string) (string, error) {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"%s/v3/accounts/%s/projects/%s/repositories/%s/",
			c.HostURL,
			strconv.Itoa(c.AccountID),
			projectID,
			repositoryID,
		),
		nil,
	)
	if err != nil {
		return "", err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return "", err
	}

	return "", err
}
