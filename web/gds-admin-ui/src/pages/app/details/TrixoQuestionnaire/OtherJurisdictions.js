const { default: Field } = require("components/Field");
const { Row, Col, Form, Button } = require("react-bootstrap");
const { useFieldArray } = require("react-hook-form");
const { isoCountries } = require("utils/country");


const OtherJurisdictions = ({ register, name, control }) => {
    const { fields, append, remove } = useFieldArray({
        name,
        control
    });


    return (
        <Row>
            <h4>Other Jurisdictions</h4>
            <p className='m-0'>
                Please add any other regulatory jurisdictions your organization complies with.
            </p>
            {
                fields.map((field, index) => (
                    <Row key={field.id} className='gy-1'>
                        <Col sm="4">
                            <Form.Group>
                                <Form.Label>National Jurisdiction</Form.Label>
                                <Field.Select type="text" register={register} name={`${name}[${index}].country`} defaultValue={field.country}>
                                    {
                                        Object.entries(isoCountries).map(([k, v]) => <option key={k} value={k}>{v}</option>)
                                    }
                                </Field.Select>
                            </Form.Group>
                        </Col>
                        <Col sm="5">
                            <Form.Group>
                                <Form.Label>Regulator Name</Form.Label>
                                <Field.Input type="text" register={register} name={`${name}[${index}].regulator_name`} defaultValue={field.regulator_name} />
                            </Form.Group>
                        </Col>
                        <Col sm="1">
                            <Form.Label className=''></Form.Label>
                            <Button onClick={() => remove(index)} style={{ marginTop: 'inherit' }} variant="danger">
                                <i className='dripicons-trash'></i>
                            </Button>
                        </Col>
                    </Row>
                ))
            }
            <Col sm="12" className='mt-2'>
                <Button onClick={() => append({
                    country: "",
                    regulator_name: ""
                })}>Add Jurisdiction</Button>
            </Col>
        </Row>
    )
}

export default OtherJurisdictions