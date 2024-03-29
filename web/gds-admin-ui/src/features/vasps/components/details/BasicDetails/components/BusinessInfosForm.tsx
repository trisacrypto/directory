import { yupResolver } from '@hookform/resolvers/yup';
import dayjs from 'dayjs';
import { Button, Col, Form, Row } from 'react-bootstrap';
import { useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import * as yup from 'yup';

import BusinessCategory from '@/components/BusinessCategory';
import Field from '@/components/Field';
import { ModalCloseButton } from '@/components/Modal';
import VASPCategory from '@/components/VaspCategory';
import { getBusinessInfosFormInitialValues } from '@/utils/form-references';
import { useUpdateVasp } from '../../../../services/update-vasp';

const schema = yup.object().shape({
    website: yup.string().url('website should be a valid url').trim().required(),
    established_on: yup.date().typeError('Date of Incorporation/Establishment should be a valid date').required(),
});

const DATE_FORMAT = 'YYYY-MM-DD';

function BusinessInfosForm({ data }: any) {
    const {
        register,
        handleSubmit,
        formState: { errors, dirtyFields },
    } = useForm({
        defaultValues: getBusinessInfosFormInitialValues(data),
        resolver: yupResolver(schema),
        mode: 'onChange',
    });
    const params = useParams<{ id: string }>();
    const { mutate: updateBusiness } = useUpdateVasp();

    const onSubmit = (values: any) => {
        if (params && params.id) {
            data && (data.established_on = dayjs(data.established_on).format(DATE_FORMAT));

            updateBusiness({
                vaspId: params.id,
                data: values,
            });
        }
    };

    return (
        <>
            <h3>Edit Business Information</h3>
            <p>
                Please enter extended and basic details about the business entity. Note that the IVMS 101 data contains
                the legal details and company name information.
            </p>
            <Form onSubmit={handleSubmit(onSubmit)}>
                <Form.Group as={Row} className="mb-2" controlId="website">
                    <Form.Label column sm="12" className="fw-normal">
                        Website
                    </Form.Label>
                    <Col sm="12">
                        <Field.Input
                            isInvalid={!!errors.website}
                            isValid={dirtyFields.website && !errors.website}
                            type="text"
                            register={register}
                            name="website"
                        />
                    </Col>
                    {errors.website ? (
                        <Form.Control.Feedback type="invalid" className="d-block">
                            {errors.website.message as string}
                        </Form.Control.Feedback>
                    ) : (
                        <Form.Text className="fst-italic">e.g: https://example.com</Form.Text>
                    )}
                </Form.Group>
                <Form.Group as={Row} className="mb-2" controlId="established_on">
                    <Form.Label column sm="12" className="fw-normal">
                        Date of Incorporation/Establishment
                    </Form.Label>
                    <Col sm="12">
                        <Field.Input
                            isInvalid={!!errors.established_on}
                            isValid={dirtyFields.established_on && !errors.established_on}
                            type="date"
                            register={register}
                            name="established_on"
                        />
                    </Col>
                    {errors.established_on && (
                        <Form.Control.Feedback type="invalid" className="d-block">
                            {errors.established_on.message as string}
                        </Form.Control.Feedback>
                    )}
                </Form.Group>
                <Form.Group as={Row} className="mb-2" controlId="business_category">
                    <Form.Label column sm="12" className="fw-normal">
                        Business Category
                    </Form.Label>
                    <Col sm="12">
                        <Field.Select register={register} name="business_category">
                            <BusinessCategory />
                        </Field.Select>
                    </Col>
                    <Form.Text>
                        Please select the entity category that most closely matches your organization.
                    </Form.Text>
                </Form.Group>
                <Form.Group as={Row} className="mb-2" controlId="vasp_categories">
                    <Form.Label column sm="12" className="fw-normal">
                        VASP Category
                    </Form.Label>
                    <Col sm="12">
                        <Field.Select htmlSize="8" register={register} multiple name="vasp_categories">
                            <VASPCategory />
                        </Field.Select>
                        <Form.Text>
                            Please select as many categories needed to represent the types of virtual asset services
                            your organization provides.
                        </Form.Text>
                    </Col>
                </Form.Group>

                <div className="mt-3 text-end">
                    <ModalCloseButton>
                        <Button variant="danger" className="me-2">
                            Cancel
                        </Button>
                    </ModalCloseButton>
                    <Button type="submit">Save</Button>
                </div>
            </Form>
        </>
    );
}

export default BusinessInfosForm;
