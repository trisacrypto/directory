import React from 'react';
import { ModalCloseButton, ModalContext } from 'components/Modal';
import OvalLoader from 'components/OvalLoader';
import {
    Button, Col, Form, Row,
} from 'react-bootstrap';
import { useForm, useWatch } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { getAdminVerificationToken, reviewVasp } from 'services/vasp';
import { useDispatch, useSelector } from 'react-redux';
import useSafeDispatch from 'hooks/useSafeDispatch';
import { reviewVaspApiResponse, reviewVaspApiResponseSuccess } from 'redux/vasp-details';
import toast from 'react-hot-toast';
import { getVaspDetailsLoadingState } from 'redux/selectors';

function ReviewForm() {
    const [verificationToken, setVerificationToken] = React.useState();
    const [errorStatusCode, setErrorStatusCode] = React.useState();
    const {
        register, handleSubmit, control, formState: { errors, dirtyFields },
    } = useForm({
        defaultValues: {
            accept: undefined,
            reject_reason: '',
        },
    });
    const params = useParams();
    const dispatch = useDispatch();
    const safeDispatch = useSafeDispatch(dispatch);
    const isLoading = useSelector(getVaspDetailsLoadingState);
    const [, setIsOpen] = React.useContext(ModalContext);
    const isMounted = React.useRef(true);
    const accept = useWatch({ control, name: 'accept' });
    const isReject = React.useMemo(() => accept === 'false', [accept]);

    React.useEffect(() => {
        if (params && params.id) {
            if (isMounted.current) {
                (async () => {
                    try {
                        const response = await getAdminVerificationToken(params.id);
                        setVerificationToken(response.data);
                    } catch (error) {
                        setErrorStatusCode(error.response.status);
                    }
                })();
            }
        }

        return () => {
            isMounted.current = false;
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [params.id]);

    const onSubmit = async (data) => {
        const payload = {
            admin_verification_token: verificationToken.admin_verification_token,
            ...data,
            accept: JSON.parse(data.accept),
        };

        safeDispatch(reviewVaspApiResponse());
        try {
            const { data } = await reviewVasp(params.id, payload);
            safeDispatch(reviewVaspApiResponseSuccess(data.status));
            toast.success(data.message, { duration: 6000 });
            setIsOpen(false);
        } catch (error) {
            console.error('[onSubmit] vasp review', error.message);
        }
    };

    if (errorStatusCode === 404) {
        return (
            <>
                <i className="mdi mdi-emoticon-sad-outline fs-2 text-center" />
                <p className="text-center">This registration is not ready for review</p>
                <div className="text-center">
                    <ModalCloseButton>
                        <Button variant="danger" className="">Close</Button>
                    </ModalCloseButton>
                </div>
            </>
        );
    }

    return !verificationToken ? <OvalLoader title="retrieving verification token..." /> : (

        <Form onSubmit={handleSubmit(onSubmit)}>
            <Row>
                <Form.Group>
                    <Form.Label>Review the registration</Form.Label>
                    <Form.Check {...register('accept')} value type="radio" label="Accept" name="accept" id="accept" className="mb-1" />
                    <Form.Check {...register('accept')} value={false} type="radio" label="Reject" name="accept" id="reject" />
                </Form.Group>
                {
                    isReject && (
                        <Form.Group className="mt-2">
                            <Form.Control
                                as="textarea"
                                rows={5}
                                isInvalid={!!errors.reject_reason}
                                {...register('reject_reason', { required: 'Please enter the rejection reason' })}
                                isValid={dirtyFields.reject_reason && !errors.reject_reason}
                                placeholder="what is the rejection reason ?"
                            />
                            {
                                errors.reject_reason ? <Form.Control.Feedback type="invalid">{errors.reject_reason.message}</Form.Control.Feedback> : <Form.Text type="invalid">* Please supply the rejection reason</Form.Text>
                            }

                        </Form.Group>
                    )
                }
            </Row>
            <Row>
                <Col>
                    <ModalCloseButton>
                        <Button variant="danger" className="mt-3 w-100">Cancel</Button>
                    </ModalCloseButton>
                </Col>
                <Col>
                    <Button type="submit" className="mt-3 w-100" disabled={isLoading}>Submit</Button>
                </Col>
            </Row>
        </Form>
    );
}

export default ReviewForm;
