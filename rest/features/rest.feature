Feature: Validate api routes

  Scenario: Check service is healthy
    When i send a GET request to "/health"
    Then the response code should be "200"

  Scenario: Create repository in github
    When i create a repository with name "dummy" and description "dummy description"
    Then the response code should be "201"

  Scenario: List user public repositories in github
    When i list all the repositories
    Then the response code should be "200"

  Scenario: Delete a user repository
    When i delete a repository with name "name"
    Then the response code should be "204"

  Scenario: List open PRs
    When i list "5" open prs from "golang" "go"
    Then the response code should be "200"
