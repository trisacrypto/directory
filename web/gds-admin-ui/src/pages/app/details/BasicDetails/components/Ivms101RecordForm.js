
import CountryOptions from 'components/CountryOptions';
import Field from 'components/Field';
import { ModalCloseButton } from 'components/Modal';
import NationalIdentifierOptions from 'components/NationalIdentifierOptions';
import React from 'react'
import { Button, Form, FormGroup } from 'react-bootstrap'
import { useForm } from 'react-hook-form';
import { getIvms101RecordInitialValues } from 'utils/form-references';
import AddressesFieldArray from './AddressesFieldArray';
import NameIdentifiersFieldArray from './NameIdentifiersFieldArray';

function Ivms101RecordForm({ data }) {
    const { register, control, handleSubmit } = useForm({
        defaultValues: getIvms101RecordInitialValues(data)
    })
    const nameIdentifiersFieldArrayRef = React.useRef()
    const localNameIdentifiersFieldArrayRef = React.useRef()
    const phoneticNameIdentifiersFieldArrayRef = React.useRef()

    const onSubmit = async (data) => {
        console.log('[onSubmit] data', data)
    }

    const handleAddLegalNamesRow = () => {
        nameIdentifiersFieldArrayRef.current.addRow()
    }

    const handleAddNewLocalNamesRow = () => {
        localNameIdentifiersFieldArrayRef.current.addRow()
    }

    const handleAddNewPhoneticNamesRow = () => {
        phoneticNameIdentifiersFieldArrayRef.current.addRow()
    }

    return (
        <div>
            <h3>Legal Person</h3>
            <p>Please enter the information that identify your organization as a Legal Person. This form represents the IVMS 101 data structure for legal persons and is strongly suggested for use as KYC information exchanged in TRISA transfers.</p>

            <Form onSubmit={handleSubmit(onSubmit)}>
                <NameIdentifiersFieldArray
                    name="name.name_identifiers"
                    register={register}
                    control={control}
                    heading={`Name Identifiers`}
                    description={`The name and type of name by which the legal person is known.`}
                    controlId={`name_identifiers`}
                    ref={nameIdentifiersFieldArrayRef}
                />

                <NameIdentifiersFieldArray
                    name="name.local_name_identifiers"
                    register={register}
                    control={control}
                    heading={`Local Name Identifiers`}
                    description={`The name by which the legal person is known using local characters.`}
                    controlId={`local_name_identifiers`}
                    ref={localNameIdentifiersFieldArrayRef}
                />

                <NameIdentifiersFieldArray
                    name="name.phonetic_name_identifiers"
                    register={register}
                    control={control}
                    heading={`Phonetic Name Identifiers`}
                    description={`A phonetic representation of the name by which the legal person is known.`}
                    controlId={`local_name_identifiers`}
                    ref={phoneticNameIdentifiersFieldArrayRef}
                />
                <div className='d-inline-flex gap-2'>
                    <Button onClick={handleAddLegalNamesRow}>Add Legal Names</Button>
                    <Button onClick={handleAddNewLocalNamesRow}>Add Local Names</Button>
                    <Button onClick={handleAddNewPhoneticNamesRow}>Add Phonetic Names</Button>
                </div>

                <AddressesFieldArray control={control} register={register} name="geographic_addresses" />

                <FormGroup>
                    <Form.Label className='fw-normal mt-3'>Customer Number</Form.Label>
                    <Field.Input register={register} name={`customer_number`} />
                    <Form.Text>TRISA specific identity number (UUID), only supplied if you're updating an existing registration request.</Form.Text>
                </FormGroup>

                <FormGroup>
                    <Form.Label className='fw-normal mt-3'>Country of Registration</Form.Label>
                    <Field.Select name={`country_of_registration`} register={register}>
                        <CountryOptions />
                    </Field.Select>
                </FormGroup>

                <FormGroup>
                    <h5 className='fw-normal mt-3'>National Identification</h5>
                    <p>Please supply a valid national identification number. TRISA recommends the use of LEI numbers. For more information, please visit <a href="https://www.gleif.org/" rel="noreferrer" target="_blank">GLEIF.org</a>
                    </p>
                    <Form.Label className='fw-normal'>Identification Number</Form.Label>
                    <Field.Input register={register} name={`national_identification.national_identifier`} />
                    <Form.Text>An identifier issued by an appropriate issuing authority.</Form.Text>
                </FormGroup>

                <FormGroup>
                    <Form.Label className='mt-3 fw-normal'>Identification Type</Form.Label>
                    <Field.Select register={register} name={`national_identification.national_identifier_type`}>
                        <NationalIdentifierOptions />
                    </Field.Select>
                </FormGroup>
                <FormGroup>
                    <Form.Label className='mt-3 fw-normal'>Country of Issue</Form.Label>
                    <Field.Select register={register} name={`national_identification.country_of_issue`}>
                        <CountryOptions />
                    </Field.Select>
                </FormGroup>
                <FormGroup>
                    <Form.Label className='fw-normal mt-3'>Registration Authority</Form.Label>
                    <Field.Input register={register} name={`national_identification.registration_authority`} />
                    <Form.Text>If the identifier is an LEI number, the ID used in the GLEIF Registration Authorities List.</Form.Text>
                </FormGroup>

                <div className='mt-3 text-end'>
                    <ModalCloseButton>
                        <Button variant='danger' className="me-2">Cancel</Button>
                    </ModalCloseButton>
                    <Button type='submit'>Save</Button>
                </div>
            </Form>
        </div>
    )
}

export default Ivms101RecordForm
