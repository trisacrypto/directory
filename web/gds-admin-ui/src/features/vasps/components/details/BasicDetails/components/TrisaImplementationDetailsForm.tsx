import { yupResolver } from '@hookform/resolvers/yup';
import React from 'react';
import { Alert, Button, Col, Form, Row } from 'react-bootstrap';
import { useForm, useWatch } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import * as yup from 'yup';

import Field from '@/components/Field';
import { ModalCloseButton } from '@/components/Modal';
import { getTrisaImplementationDetailsInitialValue } from '@/utils/form-references';
import { useUpdateVasp } from '@/features/vasps/services/update-vasp';

const trisaEndpointPattern = /^([a-zA-Z0-9.-]+):((?!(0))[0-9]+)$/;

const schema = yup.object().shape({
    trisa_endpoint: yup.string().trim().matches(trisaEndpointPattern, 'trisa endpoint is not valid'),
});

function TrisaImplementationDetailsForm({ data }: any) {
    const {
        register,
        handleSubmit,
        formState: { errors, dirtyFields, isDirty, isSubmitting },
        control,
    } = useForm({
        defaultValues: getTrisaImplementationDetailsInitialValue(data),
        resolver: yupResolver(schema),
        mode: 'all',
        reValidateMode: 'onChange',
    });
    const params = useParams<{ id: string }>();
    const commonName = useWatch({ name: 'common_name', control });
    const trisaEndpoint = useWatch({ name: 'trisa_endpoint', control });
    const [commonNameWarning, setCommonNameWarning] = React.useState<string>('');
    const { mutate: updateTrisaImplementationDetails, isError, error } = useUpdateVasp();

    React.useEffect(() => {
        const trisaEndpointUri = trisaEndpoint.split(':')[0];
        const warningMessage =
            trisaEndpointUri === commonName ? '' : 'common name should match the TRISA endpoint without the port';
        setCommonNameWarning(warningMessage);
    }, [commonName, trisaEndpoint]);

    const onSubmit = (data: any) => {
        updateTrisaImplementationDetails({
            vaspId: params.id,
            data,
        });
    };

    const handleAlertClose = () => {};

    return (
        <>
            <h3>Edit TRISA Implementation</h3>
            <p>
                Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please specify the
                details of your endpoint for certificate issuance.
            </p>
            <Form onSubmit={handleSubmit(onSubmit)}>
                <Alert variant="danger" show={isError} className="col-sm-12" onClose={handleAlertClose} dismissible>
                    {error as unknown as string}
                </Alert>
                <Form.Group as={Row} className="mb-2" controlId="trisa_endpoint">
                    <Form.Label column sm="12" className="fw-normal">
                        TRISA Endpoint
                    </Form.Label>
                    <Col sm="12">
                        <Field.Input
                            isInvalid={!!errors.trisa_endpoint}
                            isValid={dirtyFields.trisa_endpoint && !errors.trisa_endpoint}
                            type="text"
                            register={register}
                            name="trisa_endpoint"
                            placeholder="trisa.example.com:443"
                        />
                    </Col>
                    {errors.trisa_endpoint ? (
                        <Form.Control.Feedback type="invalid" className="d-block">
                            {dirtyFields.trisa_endpoint && errors.trisa_endpoint.message}
                        </Form.Control.Feedback>
                    ) : (
                        <Form.Text className="fst-italic">
                            The address and port of the TRISA endpoint for partner VASPs to connect on via gRPC.
                        </Form.Text>
                    )}
                </Form.Group>
                <Form.Group as={Row} className="mb-2" controlId="common_name">
                    <Form.Label column sm="12" className="fw-normal">
                        Certificate Common Name
                    </Form.Label>
                    <Col sm="12">
                        <Field.Input
                            isInvalid={dirtyFields.common_name && !!errors.common_name}
                            isValid={dirtyFields.common_name && !errors.common_name}
                            type="text"
                            register={register}
                            name="common_name"
                            placeholder="trisa.example.com"
                        />
                    </Col>
                    {commonNameWarning ? (
                        <Form.Text className="fst-italic text-warning">
                            <i className="dripicons-warning" /> {commonNameWarning}
                        </Form.Text>
                    ) : (
                        <Form.Text className="fst-italic">
                            The common name for the mTLS certificate. This should match the TRISA endpoint without the
                            port in most cases.
                        </Form.Text>
                    )}
                </Form.Group>
                <div className="mt-3 text-end">
                    <ModalCloseButton>
                        <Button variant="danger" className="me-2">
                            Cancel
                        </Button>
                    </ModalCloseButton>
                    <Button type="submit" disabled={!isDirty || isSubmitting}>
                        Save
                    </Button>
                </div>
            </Form>
        </>
    );
}

export default TrisaImplementationDetailsForm;
