Feature: Simple Proxy

  As an operator, I would like to monitor the number, frequency, and duration
  of connections without affecting their behavior at all.

  As a PostgreSQL user, I would like to connect without changing any of my
  settings.


  Background:
    Given postgres "pg0" is configured to trust connections
      And it is running
    Given pgtwixt is configured to connect to "pg0"
      And it is running

  Scenario: Clients connect and disconnect
     When a client connects to pgtwixt
     Then pgtwixt reports a frontend and a backend connect
     When that client disconnects
     Then pgtwixt reports a frontend and a backend disconnect

