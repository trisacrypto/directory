import React from 'react';
import Countries from '../select/Countries';
import LegalPersonName from './LegalPersonName';
import NationalIdentification from './NationalIdentification';
import AddressList from './AddressList';
import Form from 'react-bootstrap/Form';
import update from 'immutability-helper';
import { Trans } from "@lingui/macro"


const LegalPerson = ({person, onChange}) => {
  const createChangeHandler = (field) => (event) => {
    const changes = {[field]: {$set: event.target.value}};
    const updated = update(person, changes);
    onChange(null, updated);
  }

  const createNestedChangeHandler = (field) => (event, value) => {
    const changes = {[field]: {$set: value}};
    const updated = update(person, changes);
    onChange(event, updated)
  }

  return (
    <>

    <LegalPersonName
      name={person.name}
      onChange={createNestedChangeHandler('name')}
    />

    <AddressList
      addresses={person.geographic_addresses}
      onChange={createNestedChangeHandler('geographic_addresses')}
    />

    <Form.Group controlId="legalPersonCustomerNumber">
      <Form.Label><Trans>Customer Number</Trans></Form.Label>
      <Form.Control
        type="text"
        value={person.customer_number}
        onChange={createChangeHandler("customer_number")}
      />
      <Form.Text className="text-muted">
        <Trans>TRISA specific identity number (UUID), only supplied if you're updating an existing registration request.</Trans>
      </Form.Text>
    </Form.Group>

    <Form.Group controlId="legalPersonCountryOfRegistration">
        <Form.Label><Trans>Country of Registration</Trans></Form.Label>
        <Form.Control
          as="select" custom
          value={person.country_of_registration}
          onChange={createChangeHandler("country_of_registration")}
        >
          <Countries />
        </Form.Control>
      </Form.Group>

    <NationalIdentification
      data={person.national_identification}
      onChange={createNestedChangeHandler("national_identification")}
    />
    </>
  );
}

export default LegalPerson;