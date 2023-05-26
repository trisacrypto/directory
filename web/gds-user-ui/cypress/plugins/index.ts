import cucumber from 'cypress-cucumber-preprocessor';

module.exports = (on, config) => {
  on('file:preprocessor', cucumber());
};
