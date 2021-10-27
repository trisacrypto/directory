import React from 'react';
import Countries from '../select/Countries';
import AddressTypeCode from '../select/AddressTypeCode';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import update from 'immutability-helper';
import { Trans } from "@lingui/macro";


const addressLineExamples = [
  "building name/number, street name",
  "apartment or suite number",
  "city, province, postal code",
]

const Address = ({index, address, onChange, onDelete}) => {
  const createChangeHandler = (field) => (event) => {
    const changes = {[field]: {$set: event.target.value}};
    const updated = update(address, changes);
    onChange(null, updated);
  }

  const createArrayChangeHandler = (field, idx) => (event) => {
    const changes = {[field]: {[idx]: {$set: event.target.value}}};
    const updated = update(address, changes);
    onChange(null, updated);
  };

  const addressLines = address.address_line.map((line, idx) => {
    return (
      <Form.Group key={idx}>
        <Form.Control
          type="text"
          value={line}
          onChange={createArrayChangeHandler('address_line', idx)}
        />
        <Form.Text className="text-muted">
          {`Address line ${idx+1} e.g. ${addressLineExamples[idx]}`}
        </Form.Text>
      </Form.Group>
    );
  });

  return (
    <>
    <Form.Label><Trans>Address {index+1}</Trans></Form.Label>
    {addressLines}
    <Form.Row>
      <Form.Group as={Col}>
        <Form.Control
          as="select" custom
          value={address.country}
          onChange={createChangeHandler('country')}
        >
          <Countries />
        </Form.Control>
        <Form.Text className="text-muted">
          <Trans>Country</Trans>
        </Form.Text>
      </Form.Group>
      <Form.Group as={Col}>
        <Form.Control
          as="select" custom
          value={address.address_type}
          onChange={createChangeHandler('address_type')}
        >
          <AddressTypeCode />
        </Form.Control>
        <Form.Text className="text-muted">
          <Trans>Address Type</Trans>
        </Form.Text>
      </Form.Group>
      <Form.Group as={Col} xs={1}>
        <Button variant="danger" onClick={onDelete} title={`Delete Address ${index+1}`}>
          <i className="fa fa-trash"></i>
        </Button>
      </Form.Group>
    </Form.Row>
    </>
  );
}

export default Address;