import Field from 'components/Field'
import LegalPersonNameIdentifierTypeOptions from 'components/LegalPersonNameIdentifierTypeOptions'
import React from 'react'
import { Button, Col, Form, Row } from 'react-bootstrap'
import { useFieldArray } from 'react-hook-form';


const NameIdentifiersFieldArray = React.forwardRef((props, ref) => {
    const { register, control, name, controlId, description, heading } = props
    const { fields, remove, append } = useFieldArray({ control, name });

    React.useImperativeHandle(ref, () => ({
        addRow() {
            append({
                legal_person_name: "",
                legal_person_name_identifier_type: '0'
            })
        }
    }))

    return (
        <Form.Group as={Row} className="" controlId={controlId}>
            {
                fields.map((field, index) => (
                    <Row sm={12} key={field.id} className='mb-2'>
                        {
                            index < 1 && (
                                <>
                                    <Form.Label className='mb-0 fw-normal'>{heading}</Form.Label>
                                    <Form.Text className='mb-1'>{description}</Form.Text>
                                </>
                            )
                        }
                        <Col sm={8}>
                            <Field.Input type="text" register={register} name={`${name}[${index}].legal_person_name`} />
                        </Col>
                        <Col sm={3}>
                            <Field.Select register={register} name={`${name}[${index}].legal_person_name_identifier_type`}>
                                <LegalPersonNameIdentifierTypeOptions />
                            </Field.Select>
                        </Col>
                        <Col sm={1}>
                            <Button variant='danger' onClick={() => remove(index)}><i className='dripicons-trash'></i></Button>
                        </Col>
                    </Row>
                ))

            }
        </Form.Group>
    )

})

export default NameIdentifiersFieldArray