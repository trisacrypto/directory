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

declare global {
    namespace Cypress {
      interface Chainable {
        loginWith({ email, password }): Chainable<JQuery<HTMLElement>>;
      }
    }
  }