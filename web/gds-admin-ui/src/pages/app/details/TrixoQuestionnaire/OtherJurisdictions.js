import { Button, Col, Form, Row } from 'react-bootstrap';
import { useFieldArray } from 'react-hook-form';

import Field from '@/components/Field';
import { isoCountries } from '@/utils/country';

const OtherJurisdictions = ({ register, name, control }) => {
  const { fields, append, remove } = useFieldArray({
    name,
    control,
  });

  return (
    <Row>
      <h4>Other Jurisdictions</h4>
          <p className="m-0">Please add any other regulatory jurisdictions your organization complies with.</p>
          {fields.map((field, index) => (
        <Row key={field.id} className="gy-1">
          <Col sm="4">
            <Form.Group>
              <Form.Label>National Jurisdiction</Form.Label>
                          <Field.Select
                type="text"
                register={register}
                              name={`${name}[${index}].country`}
                              defaultValue={field.country}
                            >
                              {Object.entries(isoCountries).map(([k, v]) => (
                  <option key={k} value={k}>
                    {v}
                  </option>
                ))}
              </Field.Select>
            </Form.Group>
          </Col>
          <Col sm="5">
            <Form.Group>
              <Form.Label>Regulator Name</Form.Label>
                          <Field.Input
                type="text"
                register={register}
                name={`${name}[${index}].regulator_name`}
                defaultValue={field.regulator_name}
              />
            </Form.Group>
          </Col>
          <Col sm="1">
            <Form.Label className="" />
            <Button onClick={() => remove(index)} style={{ marginTop: 'inherit' }} variant="danger">
              <i className="dripicons-trash" />
            </Button>
          </Col>
        </Row>
      ))}
      <Col sm="12" className="mt-2">
        <Button
                  onClick={() =>
            append({
              country: '',
              regulator_name: '',
            })}
                >
        >
          Add Jurisdiction
        </Button>
      </Col>
    </Row>
  );
};

export default OtherJurisdictions;
