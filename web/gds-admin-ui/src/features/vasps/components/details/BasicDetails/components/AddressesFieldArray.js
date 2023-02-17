import React from 'react';
import { Button, Col, Form, FormGroup, Row } from 'react-bootstrap';
import { useFieldArray } from 'react-hook-form';

import AddressTypeOptions from '@/components/AddressTypeOptions';
import CountryOptions from '@/components/CountryOptions';
import Field from '@/components/Field';

function AddressesFieldArray({ control, name, register }) {
  const { fields, remove, append } = useFieldArray({ control, name });

  return (
    <FormGroup>
      <Form.Label className="mb-0 fw-normal mt-3">Addresses</Form.Label>
      <p className="small">Please enter at least one geographic address.</p>

      {fields.map((field, idx) => (
                    <Row sm={12} key={field.id} className="mb-2">
                        <p>Address {idx + 1}</p>
                        <Col sm={12} className="mb-2">
                            <Field.Input type="text" register={register} name={`${name}[${idx}].address_line[0]`} />
                            <Form.Text>Address line 1 e.g. building name/number, street name</Form.Text>
                        </Col>
                        <Col sm={12} className="mb-2">
                            <Field.Input type="text" register={register} name={`${name}[${idx}].address_line[1]`} />
                            <Form.Text>Address line 2 e.g. apartment or suite number</Form.Text>
                        </Col>
                        <Col sm={12} className="mb-2">
                            <Field.Input type="text" register={register} name={`${name}[${idx}].address_line[2]`} />
                            <Form.Text>Address line 3 e.g. city, province, postal code</Form.Text>
                        </Col>
                        <Col sm={6}>
                            <Field.Select name={`${name}[${idx}.country]`} register={register}>
                                <CountryOptions />
                            </Field.Select>
                            <Form.Text>Country</Form.Text>
                        </Col>
                        <Col sm={5}>
                            <Field.Select register={register} name={`${name}[${idx}].address_type`}>
                                <AddressTypeOptions />
                            </Field.Select>
                            <Form.Text>Address Type</Form.Text>
                        </Col>
                        <Col sm={1} className="ps-1">
                            <Button variant="danger" onClick={() => remove(idx)}>
                                <i className="dripicons-trash"></i>
                            </Button>
                        </Col>
                    </Row>
                ))}
      <div>
        <Button
          onClick={() =>
            append({
              address_type: 2,
              address_line: ['', '', ''],
                        country: '',
            })}
                >
        >
          Add Address
        </Button>
      </div>
    </FormGroup>
  );
}

export default AddressesFieldArray;
