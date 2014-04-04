#!/usr/bin/env ruby
# encoding: utf-8
Feature: Agreement CRUD
	In order to provide users with any value
	we need to be able to create and save an agreement

	Scenario: New agreement
		Given a new agreement
		When I save it
		Then I have an id

	Scenario: Update a specific version of an agreement
		Given an existing agreement
		When I update it
		Then an agreement is returned

	Scenario: Fetch a specific version of an agreement
		Given an existing agreement
		When I fetch its version
		Then an agreement is returned

	Scenario: Delete a specific version of an agreement
		Given an existing agreement
		When I update it
		Then an empty response is returned

	Scenario: Fetch lastest version of an agreement 
		Given an existing agreement
		When I fetch its agreement
		Then an agreement is returned

	Scenario: Fetch a users live agreements 
		Given an existing agreement
		When I fetch based on user ID
		Then at least one agreement is returned

	Scenario: Fetch an agreement versions owners 
		Given an existing agreement
		When I request an agreement versions owners
		Then I get all owners back

	Scenario: Fetch an agreements owners 
		Given an existing agreement
		When I request an agreements owners
		Then I get all owners back

	Scenario: Fetch a users completed agreements 
		Given an existing agreement
		When I update its last action as completed
		When I request a users archived agreements
		Then at least one agreement is returned