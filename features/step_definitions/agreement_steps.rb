#encoding: utf-8
require "rubygems"
require "majordomo"
require "json"

Before do
	@client = Majordomo::Client.new("tcp://localhost:5555", false)
	jsonString = ['{"title":"ruby", "freelancerID":"a0f62150-b93f-11e3-a5e2-0800200c9a66", "clientID":"a213f1c0-b95d-11e3-a5e2-0800200c9a66"}'].pack('m0')
	reply = @client.send_and_receive("Agreements", '{"Method":"POST", "Path":"/agreements/v", "Body":"'+jsonString+'"}')
	@agreement = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

After do
	@client.close
end

Given /^a new agreement$/ do

end

Given /^an existing agreement$/ do

end

When /^I save it$/ do

end

When /^I update it$/ do
	jsonString = ['{"title":"updated"}'].pack('m0')
	reply = @client.send_and_receive("Agreements", '{"Method":"PUT", "Path":"/agreements/v/'+ @agreement["versionID"]+'", "Body":"'+jsonString+'"}')
	@new_agreement = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I fetch its version$/ do
	reply = @client.send_and_receive("Agreements", '{"Method":"GET", "Path":"/agreements/v/'+ @agreement["versionID"]+'"}')
	@new_agreement = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I fetch its agreement$/ do
	reply = @client.send_and_receive("Agreements", '{"Method":"GET", "Path":"/agreements/'+ @agreement["agreementID"]+'"}')
	@new_agreement = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I fetch based on user ID$/ do
	reply = @client.send_and_receive("Agreements", '{"Method":"GET", "Path":"/user/'+ @agreement["freelancerID"]+'/agreements"}')
	@agreements = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I request an agreement versions owners$/ do
	reply = @client.send_and_receive("Agreements", '{"Method":"GET", "Path":"/agreements/v/'+ @agreement["versionID"]+'/owners"}')
	@owners = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I request an agreements owners$/ do
	reply = @client.send_and_receive("Agreements", '{"Method":"GET", "Path":"/agreements/'+ @agreement["agreementID"]+'/owners"}')
	@owners = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I delete it$/ do
	reply = @client.send_and_receive("Agreements", '{"Method":"DELETE", "Path":"/agreements/v/'+ @agreement["versionID"]+'"}')
	@new_agreement = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I update its last action as completed$/ do
	jsonString = ['{"name":"completed"}'].pack('m0')
	@client.send_and_receive("Agreements", '{"Method":"POST", "Path":"/agreements/v/'+ @agreement["versionID"]+'/action", "Body":"'+jsonString+'"}')
end

When /^I request a users archived agreements$/ do
	reply = @client.send_and_receive("Agreements", '{"Method":"GET", "Path":"/user/'+ @agreement["freelancerID"]+'/archives"}')
	@agreements = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

Then /^I have an id$/ do
	!@agreement["versionID"].empty?
end

Then /^an agreement is returned$/ do
	!@new_agreement["versionID"].nil?
end

Then /^at least one agreement is returned$/ do
	@agreements.kind_of?(Array) and @agreements.length > 0
end

Then /^an empty response is returned$/ do
	@new_agreement.empty?
end

Then /^I get all owners back$/ do
	!@owners["clientID"].empty? and !@owners["freelancerID"].empty?
end