import 'react-phone-number-input/style.css';

import { Alert, Button, Col, Fade, Form, Row } from 'react-bootstrap';
import { Controller } from 'react-hook-form';
import PhoneInput from 'react-phone-number-input';

import { ModalCloseButton } from '@/components/Modal';
import { validEmailPattern } from '@/constants';
import useContactForm from '@/features/vasps/services/use-contact-form';

const EditContactForm = ({ contactType, contact }: any) => {
    const {
        isError,
        handleAlertClose,
        onSubmit,
        handleSubmit,
        error,
        control,
        formState: { isDirty },
    } = useContactForm({ contactType, contact });

    return (
        <>
            <Fade in={isError}>
                <Alert variant="danger" show={isError} onClose={handleAlertClose} className="col-sm-12" dismissible>
                    {error as unknown as string}
                </Alert>
            </Fade>
            <Form onSubmit={handleSubmit(onSubmit)}>
                <Form.Group as={Row} className="mb-2" controlId={`${contactType}.name`}>
                    <Form.Label column sm="12" className="fw-normal">
                        Full Name
                    </Form.Label>
                    <Col sm="12">
                        <Controller
                            name="name"
                            control={control}
                            render={({ field, fieldState: { error, invalid, isDirty } }) => (
                                <>
                                    <Form.Control
                                        isInvalid={isDirty && invalid}
                                        isValid={isDirty && !invalid}
                                        type="text"
                                        placeholder="trisa"
                                        {...field}
                                    />
                                    {error ? (
                                        <Form.Control.Feedback type="invalid" className="d-block">
                                            {error.message}
                                        </Form.Control.Feedback>
                                    ) : (
                                        <Form.Text className="fst-italic">
                                            Preferred name for email communication.
                                        </Form.Text>
                                    )}
                                </>
                            )}
                        />
                    </Col>
                </Form.Group>
                <Form.Group as={Row} className="mb-2" controlId="email">
                    <Form.Label column sm="12" className="fw-normal">
                        Email address
                    </Form.Label>
                    <Col sm="12">
                        <Controller
                            name="email"
                            control={control}
                            rules={{
                                required: 'Email is required',
                                pattern: {
                                    value: validEmailPattern,
                                    message: 'Please enter a valid email',
                                },
                            }}
                            render={({ field, fieldState: { invalid, error, isDirty } }) => (
                                <>
                                    <Form.Control
                                        isInvalid={isDirty && invalid}
                                        isValid={isDirty && !invalid}
                                        type="email"
                                        placeholder="trisa@mail.com"
                                        {...field}
                                    />
                                    {error ? (
                                        <Form.Control.Feedback type="invalid" className="d-block">
                                            {error.message}
                                        </Form.Control.Feedback>
                                    ) : (
                                        <Form.Text className="fst-italic">
                                            Please use the email address associated with your organization.
                                        </Form.Text>
                                    )}
                                </>
                            )}
                        />
                    </Col>
                </Form.Group>
                <Form.Group as={Row} className="mb-2" controlId="phone">
                    <Form.Label column sm="12" className="fw-normal">
                        Phone Number
                    </Form.Label>
                    <Col sm="12">
                        <Controller
                            control={control}
                            name="phone"
                            render={({ field, fieldState: { invalid, isDirty } }) => (
                                <PhoneInput
                                    inputComponent={Form.Control}
                                    isInvalid={isDirty && invalid}
                                    isValid={isDirty && !invalid}
                                    type="tel"
                                    international={true} // set the international number format
                                    limitMaxLength={true}
                                    defaultCountry="US"
                                    {...field}
                                />
                            )}
                        />
                    </Col>
                </Form.Group>
                <div className="text-end">
                    <ModalCloseButton>
                        <Button variant="danger">Cancel</Button>
                    </ModalCloseButton>
                    <Button type="submit" variant="primary" className="ms-2" disabled={!isDirty}>
                        Submit
                    </Button>
                </div>
            </Form>
        </>
    );
};

export default EditContactForm;
