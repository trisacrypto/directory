Feature: Registration Form

I want to navigate and complete the registration form

Scenario: Registration Form

Given I'm logged in
When I click to start or complete the registration process
Then I should see the registration form
When I click the Legal Person stepper label
Then I should see the Legal Person form
When I click the Contacts stepper label
Then I should see the Contacts form
When I click the TRISA stepper label
Then I should see the TRISA form
When I click the TRIXO stepper label
Then I should see the Trixo form
When I click the Review stepper label
Then I should see the Review page
When I click the Basic Details stepper label
Then I should see the Basic Details form
When I click Save & Next at the bottom of the Basic Details form
Then I should see the Legal Person form
When I click Save & Previous at the bottom of the Legal Person form
Then I should see the Basic Details form
When I type info into the Basic Details form
And I click the Contacts stepper label
Then I should see the Unsave Changes Modal
And I should stay on the Basic Details page if I click Cancel
And I should see the Contacts form if I click Continue
# When I navigate to the Review page
