import { Given, When, Then } from "@badeball/cypress-cucumber-preprocessor";

Given("I'm logged in", () => {
    cy.login();
    cy.location("pathname").should("eq", "/dashboard/overview");
});

When("I click to start or complete the registration process", () => {
    cy.get('[data-cy="needs-attention"]').within(() => {
        cy.get('a').click();
    });
});

Then("I should see the registration form", () => {
    cy.location("pathname").should("eq", "/dashboard/certificate/registration");
});

When("I click the Legal Person stepper label", () => {
    cy.get('[data-cy="step-2-bttn"]').click();
});

Then("I should see the Legal Person form", () => {
    cy.get('[data-cy="legal-person-form"]').should("exist");
});

When("I click the Contacts stepper label", () => {
    cy.get('[data-cy="step-3-bttn"]').click();
});

Then("I should see the Contacts form", () => {
    cy.get('[data-cy="contacts-form"]').should("exist");
});

When("I click the TRISA stepper label", () => {
    cy.get('[data-cy="step-4-bttn"]').click();
});

Then("I should see the TRISA form", () => {
    cy.get('[data-cy="trisa-form"]').should("exist");
});

When("I click the TRIXO stepper label", () => {
    cy.get('[data-cy="step-5-bttn"]').click();
});

Then("I should see the Trixo form", () => {
    cy.get('[data-cy="trixo-form"]').should("exist");
});

When("I click the Review stepper label", () => {
    cy.get('[data-cy=step-6-bttn]').click();
});

Then("I should see the Review page", () => {
    cy.get('[data-cy="review-page"]').should("exist");
});

When("I click the Basic Details stepper label", () => {
    cy.get('[data-cy="step-1-bttn"]').click({ force: true });
});

Then("I should see the Basic Details form", () => {
    cy.get('[data-cy="basic-details-form"]').should("exist");
});

When("I type info into the Basic Details form", () => {
    cy.get('input[name="organization_name"]').type('Test Company');
});

Then("I click to navigate to the Contacts form without saving changes", () => {
    cy.get('[data-cy="step-3-bttn"]').click();
});

Then("I should see the Unsaved changes alert modal", () => {
    cy.contains("Unsaved changes alert").should("exist");
});

Then("I should stay on the Basic Details page if I click Cancel", () => {
    cy.get('[data-cy="cancel-bttn"]').click({ force: true });
    cy.get('[data-cy="basic-details-form"]').should("exist");
});

Then("I should see the Contacts form if I click Continue", () => {
    cy.get('[data-cy="step-3-bttn"]').click();
    cy.contains("Unsaved changes alert").should("exist");
    cy.get('[data-cy="continue-bttn"]').click({ force: true });
    cy.get('[data-cy="contacts-form"]').should("exist");
});

When("I complete a field in the Contacts Form", () => {
    cy.get('input[name="contacts.legal.name"]').type('Kamala Khan');
});

Then("I save changes to the form", () => {
    cy.contains("Save & Next").click();
});

When("I click the Clear & Reset section button on the Contacts form", () => {
    cy.get('[data-cy="trisa-form"]').should("exist");
    cy.contains("Save & Previous").click();
    cy.get('[data-cy="contacts-form"]').should("exist");
    cy.contains("Clear & Reset Section").click({ force: true });
    cy.contains("Reset").click({ force: true });
});

Then("I should not see any data in the Contacts form", () => {
    cy.get('input[name="contacts.legal.name"]').should("be.empty");
});

When("I complete the required fields in the Basic Details Form", () => {
    cy.get('[data-cy=step-1-bttn]').click({ force: true });
    cy.get('input[name="organization_name"]').type('Test Company');
    cy.get('input[name="website"]').type('https://www.test.com');
    cy.get('input[name="established_on"]').type('2023-01-01');
    cy.contains("Save & Next").click();
});

When("I complete the required fields in the Legal Form", () => {
    cy.get('[data-cy="basic-details-form"]').should("exist");
    cy.get('input[name="entity.geographic_addresses[0].address_line[0]"]').type('123 Test Street');
    cy.get('input[name="entity.geographic_addresses[0].town_name"]').type('La Ciudad');
    cy.get('input[name="entity.geographic_addresses[0].country_sub_division"]').type('MA');
    cy.contains("Select").click().type("United States{enter}");
    cy.get('input[name="entity.geographic_addresses[0].post_code"]').type('12345');
    cy.contains("Select a country").click().type("United States{enter}");
    cy.get('input[name="entity.national_identification.national_identifier"]').type('123456789');
    cy.contains("Save & Next").click();
});

When("I complete the required fields in the Contacts Form", () => {
    cy.get('[data-cy="contacts-form"]').should("exist");
    cy.get('input[name="contacts.legal.name"]').type('Kamala Khan');
    cy.get('input[name="contacts.legal.email"]').type('kamala@test.com');
    cy.get('input[name="contacts.legal.phone"]').type('555555555555');
    cy.get('input[name="contacts.technical.name"]').type('Bruno Carrelli');
    cy.get('input[name="contacts.technical.email"]').type('bruno@test.com');
    cy.get('input[name="contacts.technical.phone"]').type('555555555555');
    cy.contains("Save & Next").click();
});

When("I complete the required fields in the TRISA Form", () => {
    cy.get('[data-cy="trisa-form"]').should("exist");
    cy.get('input[name="testnet.endpoint"]').type('test.name:4477');
    cy.get('input[name="testnet.common_name"]').type('test.name');
    cy.contains("Save & Next").click();
});

When("I complete the required fields in the TRIXO Form", () => {
    cy.get('[data-cy="trixo-form"]').should("exist");
    cy.contains("Select").click();
    cy.contains("United States").click();
    cy.get('input[name="trixo.primary_regulator"]').type('Kamala Khan');
    cy.contains("Save & Next").click();
});

When("I click Save & Next on the Review page", () => {
    cy.get('[data-cy="review-page"]').should("exist");
    cy.contains("Save & Next").click();
});

Then("I should see the Registration Submission page", () => {
    cy.contains('Registration Submission').should('exist');
});

When("I click the Back to Review Section button", () => {
    cy.get('[data-cy="back-to-review-section"]').click();
});

Then("I should be returned to the Review page", () => {
    cy.get('[data-cy="review-page"]').should("exist");
});

When("I click the Clear & Reset form button", () => {
    cy.contains("Clear & Reset Form").click({ force: true });
    cy.contains("Reset").click({ force: true });
});

Then("I should not see any data in the Basic Details form", () => {
    cy.get('[data-cy="step-1-bttn"]').click({ force: true });
    cy.get('input[name="organization_name"]').should('be.empty');
    cy.get('input[name="website"]').should('be.empty');
    cy.get('input[name="established_on"]').should('be.empty');
});
