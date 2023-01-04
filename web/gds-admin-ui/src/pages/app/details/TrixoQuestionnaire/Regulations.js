import { Button, Col, Form, Row } from 'react-bootstrap';
import { useFieldArray } from 'react-hook-form';

import Field from '@/components/Field';

const Regulations = ({ register, name, control }) => {
  const { fields, append, remove } = useFieldArray({
    name,
    control,
  });

  return (
    <Form.Group as={Row} className="mb-3">
      <Form.Label column sm="12" className="fw-normal">
        Applicable Regulations
      </Form.Label>
      <Form.Text className="mb-1">
        Please specify the applicable regulation(s) for Travel Rule standards compliance, e.g. "FATF
        Recommendation 16"
      </Form.Text>
      {fields.map((field, index) => (
        <Row key={field.id} className="gy-1">
          <Col sm="9">
            <Field.Input
              type="text"
              register={register}
              name={`${name}[${index}].name`}
              defaultValue={field.name}
            />
          </Col>
          <Col sm="1">
            <Button variant="danger" onClick={() => remove(index)}>
              <i className="dripicons-trash" />
            </Button>
          </Col>
        </Row>
      ))}
      <Col sm="12" className="mt-2">
        <Button onClick={() => append({ name: '' })}>Add Regulation</Button>
      </Col>
    </Form.Group>
  );
};

export default Regulations;
