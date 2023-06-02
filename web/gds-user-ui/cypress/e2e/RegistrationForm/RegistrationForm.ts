import { Given, When, Then } from "@badeball/cypress-cucumber-preprocessor"

Given("I'm logged in", () => {
    cy.login()
    cy.location("pathname").should("eq", "/dashboard/overview")
});

When("I click to start or complete the registration process", () => {
    cy.get("[data-cy=needs-attention]").within(() => {
        cy.get('a').click()
    });
});

Then("I should see the registration form", () => {
    cy.location("pathname").should("eq", "/dashboard/certificate/registration")
});

When("I click the Legal Person stepper label", () => {

});

Then("I should see the Legal Person form", () => {
    cy.get("[data-cy=legal-person-form]").should("exist")
});