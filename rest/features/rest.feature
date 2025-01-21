Feature: Validate api routes

  Scenario: Check service is healthy
    When i send a GET request to "/health"
    Then the response code should be "200"
