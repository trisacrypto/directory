import { Form, Row, Col, Button, Alert, Fade } from "react-bootstrap"
import { Controller, useForm } from "react-hook-form"
import 'react-phone-number-input/style.css'
import PhoneInput from 'react-phone-number-input'
import { ModalCloseButton, ModalContext } from "components/Modal"
import { useDispatch, useSelector } from "react-redux"
import { getContactErrorState, getContacts, getVaspDetailsLoadingState } from "redux/selectors"
import { getContactInitialValues } from "utils/form-references"
import { validEmailPattern } from "constants/index"
import useSafeDispatch from "hooks/useSafeDispatch"
import React from "react"
import { clearContactErrorMessage, updateContact } from "redux/vasp-details"
import { useParams } from "react-router-dom"

// eslint-disable-next-line no-useless-escape

const EditContactForm = ({ contactType }) => {
    const contacts = useSelector(getContacts)
    const { control, handleSubmit, formState: { isDirty } } = useForm({
        defaultValues: getContactInitialValues(contacts[contactType]),
        mode: 'onChange'
    })
    const params = useParams()
    const dispatch = useDispatch()
    const safeDispatch = useSafeDispatch(dispatch)
    const [isOpen, setIsOpen] = React.useContext(ModalContext)
    const isLoading = useSelector(getVaspDetailsLoadingState)
    const getVaspDetailsError = useSelector(getContactErrorState)

    React.useEffect(() => {

        return () => {
            dispatch(clearContactErrorMessage())
        }
    }, [])

    const onSubmit = (data) => {
        if (params && params.id) {
            safeDispatch(updateContact({
                vaspId: params.id,
                contactType,
                data,
                setIsOpen
            }))
        }
    }

    const handleAlertClose = () => {
        dispatch(clearContactErrorMessage())
    }

    return (
        <>
            <Fade in={!!getVaspDetailsError?.error}>
                <Alert variant="danger" show={!!getVaspDetailsError?.error} onClose={handleAlertClose} className='col-sm-12' dismissible>
                    {getVaspDetailsError?.error}
                </Alert>
            </Fade>
            <Form onSubmit={handleSubmit(onSubmit)}>
                <Form.Group as={Row} className="mb-2" controlId={`${contactType}.name`}>
                    <Form.Label column sm="12" className="fw-normal">Full Name</Form.Label>
                    <Col sm="12">
                        <Controller
                            name="name"
                            control={control}
                            render={({ field, fieldState: { error, invalid, isDirty } }) => {
                                return (
                                    <>
                                        <Form.Control
                                            isInvalid={isDirty && invalid}
                                            isValid={isDirty && !invalid}
                                            type="text"
                                            placeholder="trisa"
                                            {...field}
                                        />
                                        {error ?
                                            <Form.Control.Feedback type='invalid' className='d-block'>{error.message}</Form.Control.Feedback> :
                                            <Form.Text className='fst-italic'>Preferred name for email communication.
                                            </Form.Text>}
                                    </>
                                )
                            }}
                        />
                    </Col>
                </Form.Group>
                <Form.Group as={Row} className="mb-2" controlId='email'>
                    <Form.Label column sm="12" className="fw-normal">Email address</Form.Label>
                    <Col sm="12">
                        <Controller
                            name="email"
                            control={control}
                            rules={{
                                required: 'Email is required',
                                pattern: {
                                    value: validEmailPattern,
                                    message: 'Please enter a valid email',
                                }
                            }}
                            render={({ field, fieldState: { invalid, error, isDirty } }) => {
                                return (
                                    <>
                                        <Form.Control
                                            isInvalid={isDirty && invalid}
                                            isValid={isDirty && !invalid}
                                            type="email"
                                            placeholder="trisa@mail.com"
                                            {...field}
                                        />
                                        {error ?
                                            <Form.Control.Feedback type='invalid' className='d-block'>{error.message}</Form.Control.Feedback> :
                                            <Form.Text className='fst-italic'>Please use the email address associated with your organization.
                                            </Form.Text>}
                                    </>
                                )
                            }}
                        />
                    </Col>
                </Form.Group>
                <Form.Group as={Row} className="mb-2" controlId='phone'>
                    <Form.Label column sm="12" className="fw-normal">Phone Number</Form.Label>
                    <Col sm="12">
                        <Controller
                            control={control}
                            name="phone"
                            render={({ field, fieldState: { invalid, isDirty } }) => (
                                <>
                                    <PhoneInput
                                        inputComponent={Form.Control}
                                        isInvalid={isDirty && invalid}
                                        isValid={isDirty && !invalid}
                                        type="tel"
                                        name="phone"
                                        international={true} // set the international number format
                                        limitMaxLength={true}
                                        defaultCountry='US'
                                        {...field}
                                    />
                                </>
                            )}
                        />
                    </Col>
                </Form.Group>
                <div className='text-end'>
                    <ModalCloseButton>
                        <Button variant='danger'>Cancel</Button>
                    </ModalCloseButton>
                    <Button type='submit' variant='primary' className='ms-2' disabled={!isDirty || isLoading}>Submit</Button>
                </div>
            </Form>
        </>
    )
}

export default EditContactForm
