/// <reference types="cypress" />
// ***********************************************

import '@testing-library/cypress/add-commands';

// loginWith is a command that may be used to log in a user with a given email and password.
Cypress.Commands.add('loginWith', ({ email, password }) =>
    cy.visit('/')
      .get('[data-cy="nav-login-bttn"]').click().location('pathname').should('eq', '/auth/login')
      .get('[data-cy="email"]').type(email)
      .get('[data-cy="password"]').type(password)
      .get('[data-cy="login-btn"]').click()
);

Cypress.Commands.add('login', () => {
  cy.visit('/')
    .get('[data-cy="nav-login-bttn"]').click().location('pathname').should('eq', '/auth/login')
    cy.fixture('user.json').then((user) => {
      cy.get('[data-cy="email"]').type(user.email)
      .get('[data-cy="password"]').type(user.password)
      .get('[data-cy="login-btn"]').click()
    });
  });

declare global {
    namespace Cypress {
      interface Chainable {
        loginWith({ email, password }): Chainable<JQuery<HTMLElement>>;
        login(): Chainable<JQuery<HTMLElement>>;
      }
    }
  }