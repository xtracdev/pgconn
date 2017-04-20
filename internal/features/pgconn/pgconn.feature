Feature: Connect to postgres
  Scenario:
    Given a running postgres instance
    When provide a connection string for the running instance
    Then a connection is returned
    And I can select the system timestamp

  Scenario:
    Given a connection string with no listener
    When I connect to no listener
    Then an error is returned

  Scenario:
    Given a loss of database connectivity
    When I detect I've lost connectivity
    Then I can reconnect
    And I can select data after reconnecting