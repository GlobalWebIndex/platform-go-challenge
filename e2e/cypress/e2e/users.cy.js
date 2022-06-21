/// <reference types="cypress"/>

context('User', () => {
  it('Add user - POST', () => {
    cy.fixture('users').as('usersJSON').then(dataUsers => {
      for (const userData of dataUsers.correctUsers) {
        cy.request('POST', '/users/', userData).as('addUserRequest');
        cy.get('@addUserRequest').then(resp => {
          expect(resp.status).to.eq(200);
          console.log(resp)
        });
      }
    })
  });
  it('Login user - POST', () => {
    cy.fixture('users').as('usersJSON').then(dataUsers => {
      for (const userData of dataUsers.correctUsers) {
        cy.request('POST', '/users/', userData).as('loginRequest');
        cy.get('@loginRequest').then(resp => {
          expect(resp.status).to.eq(200);
          console.log(resp)
        });
      }
    });
  })
})