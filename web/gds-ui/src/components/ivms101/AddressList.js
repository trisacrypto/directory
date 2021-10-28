import React from 'react';
import Address from './Address';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import update from 'immutability-helper';
import { Trans } from "@lingui/macro";


const AddressList = ({addresses, onChange}) => {
  const createArrayChangeHandler = (idx) => (event, value) => {
    const changes = {[idx]: {$set: value}};
    const updated = update(addresses, changes);
    onChange(event, updated)
  };

  const createArrayRemoveHandler = (idx) => (event) => {
    const changes = {$splice: [[idx, 1]]};
    const updated = update(addresses, changes);
    onChange(null, updated);
  };

  const addAddressHandler = (event) => {
    const changes = {$push: [{address_type: 2, address_line: ["", "", ""], country: ""}]};
    const updated = update(addresses, changes);
    onChange(null, updated);
  };

  const renderedAddresses = addresses.map((address, idx) => {
    return (
      <Address
        key={idx}
        index={idx}
        address={address}
        onChange={createArrayChangeHandler(idx)}
        onDelete={createArrayRemoveHandler(idx)}
      />
    );
  });

  return (
    <fieldset>
      <legend className="subsublegend"><Trans>Addresses</Trans></legend>
      <p>
        <Trans>Please enter at least one geographic address.</Trans>
      </p>
      {renderedAddresses}
      <Form.Group>
        <Button size="sm" variant="primary" onClick={addAddressHandler}><Trans>Add Address</Trans></Button>
      </Form.Group>
    </fieldset>
  );
}

export default AddressList;