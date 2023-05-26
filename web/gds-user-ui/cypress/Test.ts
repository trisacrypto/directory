import { Given, Then } from "cypress-cucumber-preprocessor/steps";

Given("I navigate to homepage", () => {
    cy.visit('http://localhost:3000')
})

Then("I should expect to see the TRISA Global Directory Service", () => {
    cy.get('h2').should('contain', 'TRISA Global Directory Service')
})