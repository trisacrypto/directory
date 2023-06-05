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
When I type info into the Basic Details form
Then I click to navigate to the Contacts form without saving changes
Then I should see the Unsaved changes alert modal
Then I should stay on the Basic Details page if I click Cancel
Then I should see the Contacts form if I click Continue

When I complete some fields in the Contacts Form
Then I save changes to the form
When I click the Clear & Reset section button on the Contacts form
Then I should not see any data in the Contacts form

When I complete the required fields in the Basic Details Form
When I complete the required fields in the Legal Form
When I complete the required fields in the Contacts Form
When I complete the required fields in the TRISA Form
When I complete the required fields in the TRIXO Form
When I click Save & Next on the Review page
Then I should see the Registration Submission page
When I click the Back to Review Section button
Then I should be returned to the Review page
When I click the Clear & Reset form button
Then I should be taken to the Basic Details form
Then I should not see any data in the Basic Details form
