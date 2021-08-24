import React from 'react';
import update from 'immutability-helper';
import Form from 'react-bootstrap/Form';
import Countries from '../select/Countries';
import NationalIdentifierTypeCode from '../select/NationalIdentifierTypeCode';


const NationalIdentification = ({data, onChange}) => {
  const createChangeHandler = (field) => (event) => {
    const changes = {[field]: {$set: event.target.value}};
    const updated = update(data, changes);
    onChange(null, updated);
  }

  return (
    <fieldset>
      <legend className="subsublegend">National Identification</legend>
      <p>
        Please supply a valid national identification number. TRISA recommends the use of
        LEI numbers. For more information, please visit <a href="https://www.gleif.org/" rel="noreferrer" target="_blank">GLEIF.org</a>.
      </p>
      <Form.Group>
        <Form.Label>Identification Number</Form.Label>
        <Form.Control
          type="text"
          value={data.national_identifier}
          onChange={createChangeHandler('national_identifier')}
        />
        <Form.Text className="text-muted">
          An identifier issued by an appropriate issuing authority.
        </Form.Text>
      </Form.Group>
      <Form.Group>
        <Form.Label>Identification Type</Form.Label>
        <Form.Control
          as="select" custom
          value={data.national_identifier_type}
          onChange={createChangeHandler('national_identifier_type')}
        >
          <NationalIdentifierTypeCode />
        </Form.Control>
      </Form.Group>
      <Form.Group>
        <Form.Label>Country of Issue</Form.Label>
        <Form.Control
          as="select" custom
          value={data.country_of_issue}
          onChange={createChangeHandler('country_of_issue')}
        >
          <Countries />
        </Form.Control>
      </Form.Group>
      <Form.Group>
        <Form.Label>Registration Authority</Form.Label>
        <Form.Control
          type="text"
          value={data.registration_authority}
          onChange={createChangeHandler('registration_authority')}
        />
        <Form.Text className="text-muted">
          If the identifier is an LEI number, the ID used in the GLEIF Registration Authorities List.
        </Form.Text>
      </Form.Group>
    </fieldset>
  );
}

export default NationalIdentification;