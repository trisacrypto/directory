import { yupResolver } from '@hookform/resolvers/yup';
import React from 'react';
import { Button, Form, Modal, OverlayTrigger, Tooltip } from 'react-bootstrap';
import { useForm, useWatch } from 'react-hook-form';
import toast from 'react-hot-toast';
import * as yup from 'yup';

import { useModal } from '@/contexts/modal';
import { APICore } from '@/helpers/api/apiCore';

import { getCookie } from '../utils';
import { captureException } from '@sentry/react';

const deliverCertsLabel = (
    <>
        Redeliver certificates
        <OverlayTrigger
            placement="right"
            overlay={<Tooltip>Sends PKCS12 encrypted certs to technical contact if still available</Tooltip>}>
            <span className="d-inline-block mdi mdi-help-circle-outline" style={{ marginLeft: '.3rem' }} />
        </OverlayTrigger>
    </>
);

const rejectionNoticeLabel = (
    <>
        Rejection notice
        <OverlayTrigger
            placement="right"
            overlay={<Tooltip>Sends a rejection notice to all verified contacts</Tooltip>}>
            <span className="d-inline-block mdi mdi-help-circle-outline" style={{ marginLeft: '.3rem' }} />
        </OverlayTrigger>
    </>
);

const reviewLabel = (
    <>
        Review request
        <OverlayTrigger placement="right" overlay={<Tooltip>Sends registration request to TRISA admins</Tooltip>}>
            <span className="d-inline-block mdi mdi-help-circle-outline" style={{ marginLeft: '.3rem' }} />
        </OverlayTrigger>
    </>
);

const verifyContactLabel = (
    <>
        Verify contacts
        <OverlayTrigger
            placement="right"
            overlay={<Tooltip>Sends verification link to all unverified contacts</Tooltip>}>
            <span className="d-inline-block mdi mdi-help-circle-outline" style={{ marginLeft: '.3rem' }} />
        </OverlayTrigger>
    </>
);

const schemaResolver = yupResolver(
    yup.object().shape({
        rejection_reason: yup
            .string()
            .trim()
            .test('test', 'Please enter the rejection email reason.', (value, context) => {
                if (context.parent.email_type === 'rejection' && !value) {
                    return false;
                }
                return true;
            }),
    })
);

const api = new APICore();

function ResendEmail() {
    const [isSubmitting, setIsSubmitting] = React.useState(false);
    const {
        register,
        handleSubmit,
        formState: { errors },
        control,
    } = useForm({
        shouldUnregister: true,
        resolver: schemaResolver,
        defaultValues: {
            rejection_reason: '',
            email_type: null,
        },
    });
    const { toggle, vasp, closeSendEmailModal } = useModal();
    const emailType = useWatch({
        control,
        name: 'email_type',
    });

    const onSubmit = (data: any) => {
        const cookie = getCookie('csrf_token');
        const payload = {
            action: data.email_type,
            reason: '',
        };
        setIsSubmitting(true);
        api.create(`/vasps/${vasp?.id}/resend`, payload, {
            headers: {
                'X-CSRF-TOKEN': cookie,
            },
        })
            .then((res) => {
                setIsSubmitting(false);
                toast.success('Email sent successfully');
            })
            .catch((err) => {
                setIsSubmitting(false);
                captureException(err);
            });
    };

    return (
        <div>
            <Modal show={toggle} dialogClassName="modal-right">
                <Form noValidate onSubmit={handleSubmit(onSubmit)}>
                    <Modal.Header onHide={closeSendEmailModal} closeButton>
                        <h4 className="modal-title">Resend Email</h4>
                    </Modal.Header>
                    <Modal.Body>
                        <h5 className="mb-3">{vasp?.name}</h5>
                        <p>Select admin email(s) to resend:</p>
                        <div className="">
                            <div>
                                <Form.Label htmlFor="verifyContact" />
                                <Form.Check
                                    {...register('email_type')}
                                    id="verifyContact"
                                    type="radio"
                                    value="verify_contact"
                                    label={verifyContactLabel}
                                    required
                                />

                                <Form.Label htmlFor="review" />
                                <Form.Check
                                    {...register('email_type')}
                                    id="review"
                                    type="radio"
                                    value="review"
                                    label={reviewLabel}
                                    required
                                />

                                <Form.Label htmlFor="certificateDeliver" />
                                <Form.Check
                                    {...register('email_type')}
                                    id="certificateDeliver"
                                    type="radio"
                                    label={deliverCertsLabel}
                                    value="deliver_certs"
                                    required
                                />

                                <Form.Label htmlFor="rejection" />
                                <Form.Check
                                    {...register('email_type')}
                                    id="rejection"
                                    type="radio"
                                    className="mb-1"
                                    value="rejection"
                                    label={rejectionNoticeLabel}
                                    required
                                />
                            </div>
                            {emailType === 'rejection' && (
                                <Form.Group className="mb-3" controlId="validationRejectionReason">
                                    <Form.Control
                                        {...register('rejection_reason')}
                                        as="textarea"
                                        placeholder="Rejection email reason"
                                        rows={5}
                                        isInvalid={!!(errors && errors.rejection_reason)}
                                    />
                                    {errors && errors.rejection_reason ? (
                                        <Form.Control.Feedback type="invalid">
                                            {errors.rejection_reason.message}
                                        </Form.Control.Feedback>
                                    ) : (
                                        <Form.Text id="rejectionHelpBlock" muted>
                                            * Please supply a rejection reason to resend the email.
                                        </Form.Text>
                                    )}
                                </Form.Group>
                            )}
                        </div>
                    </Modal.Body>
                    <Modal.Footer>
                        <Button variant="light" onClick={closeSendEmailModal}>
                            Cancel
                        </Button>{' '}
                        <Button type="submit" variant="primary" disabled={isSubmitting}>
                            {isSubmitting ? 'Submitting...' : 'Submit'}
                        </Button>
                    </Modal.Footer>
                </Form>
            </Modal>
        </div>
    );
}

export default ResendEmail;
