import { Button, Col, Form, Row, Alert } from 'react-bootstrap';
import { useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';

import Currencies from '@/components/Currencies';
import Field from '@/components/Field';
import { ModalCloseButton } from '@/components/Modal';
import { isoCountries } from '@/utils/country';
import { getTrixoFormInitialValues } from '@/utils/form-references';

import OtherJurisdictions from './OtherJurisdictions';
import Regulations from './Regulations';
import { useUpdateVasp } from '@/features/vasps/services/update-vasp';

const flatArray = (data: any) => data && data.map((d: any) => d.name);

function TrixoEditForm({ data }: any) {
    const params = useParams<{ id: string }>();
    const {
        register,
        control,
        handleSubmit,
        formState: { isDirty },
    } = useForm({
        defaultValues: getTrixoFormInitialValues(data),
    });

    const { mutate: updateTrixoForm, isLoading, error, isError } = useUpdateVasp();

    const onSubmit = async (data: any) => {
        const payload = {
            ...data,
            applicable_regulations: flatArray(data.applicable_regulations),
            compliance_threshold: parseFloat(data.compliance_threshold),
            kyc_threshold: parseFloat(data.kyc_threshold),
        };

        updateTrixoForm({
            vaspId: params.id,
            data: { trixo: payload },
        });
    };

    return (
        <Form onSubmit={handleSubmit(onSubmit)}>
            <Alert variant="danger" show={isError} className="col-sm-12" dismissible>
                {error as unknown as string}
            </Alert>
            <Form.Group as={Row} className="mb-3" controlId="primary_national_jurisdiction">
                <Form.Label column sm="12" className="fw-normal">
                    Primary National Jurisdiction
                </Form.Label>
                <Col sm="10">
                    <Field.Select register={register} name="primary_national_jurisdiction">
                        <option />
                        {Object.entries(isoCountries).map(([k, v]) => (
                            <option key={k} value={k}>
                                {v}
                            </option>
                        ))}
                    </Field.Select>
                </Col>
            </Form.Group>
            <Form.Group as={Row} className="mb-3" controlId="primary_regulator">
                <Form.Label column sm="12" className="fw-normal">
                    Name of Primary Regulator
                </Form.Label>
                <Col sm="10">
                    <Field.Input type="text" register={register} name="primary_regulator" />
                </Col>
                <Form.Text>
                    The name of primary regulator or supervisory authority for your national jurisdiction
                </Form.Text>
            </Form.Group>
            <OtherJurisdictions register={register} name="other_jurisdictions" control={control} />
            <Form.Group as={Row} className="mb-3 mt-2">
                <Form.Label column sm="12" className="fw-normal">
                    Is your organization permitted to send and/or receive transfers of virtual assets in the
                    jurisdictions in which it operates?
                </Form.Label>
                <Col sm="10">
                    <Field.Select register={register} name="financial_transfers_permitted">
                        <option value=""></option>
                        <option value="yes">Yes</option>
                        <option value="partial">Partially</option>
                        <option value="no">No</option>
                    </Field.Select>
                </Col>
            </Form.Group>

            <Form.Group as={Row} className="mb-3">
                <h4>CDD & Travel Rule Policies</h4>
                <Form.Label column sm="12" className="fw-normal">
                    Does your organization have a program that sets minimum AML, CFT, KYC/CDD and Sanctions standards
                    per the requirements of the jurisdiction(s) regulatory regimes where it is
                    licensed/approved/registered?
                </Form.Label>
                <Col sm="10">
                    <Field.Select register={register} name="has_required_regulatory_program">
                        <option value=""></option>
                        <option value="yes">Yes</option>
                        <option value="partial">Partially Implemented</option>
                        <option value="no">No</option>
                    </Field.Select>
                </Col>
            </Form.Group>

            <Form.Group as={Row} className="mb-3" controlId="conductsCustomerKYC">
                <Form.Label column sm="12" className="fw-normal">
                    Does your organization conduct KYC/CDD before permitting its customers to send/receive virtual asset
                    transfers?
                </Form.Label>
                <Col>
                    <Field.Switch
                        type="switch"
                        label="Conducts KYC before virtual asset transfers"
                        register={register}
                        name="conducts_customer_kyc"
                    />
                </Col>
            </Form.Group>

            <Form.Group as={Row} className="mb-3">
                <Form.Label column sm="12" className="fw-normal">
                    At what threshold and currency does your organization conduct KYC?
                </Form.Label>
                <Col sm="7">
                    <Field.Input type="number" register={register} name="kyc_threshold" />
                </Col>
                <Col sm="3">
                    <Field.Select register={register} name="kyc_threshold_currency">
                        <Currencies />
                    </Field.Select>
                </Col>
                <Form.Text>
                    Threshold to conduct KYC before permitting the customer to send/receive virtual asset transfers
                </Form.Text>
            </Form.Group>

            <Form.Group as={Row} className="mb-3" controlId="mustComplyTravelRule">
                <Form.Label column sm="12" className="fw-normal">
                    Is your organization required to comply with the application of the Travel Rule standards in the
                    jurisdiction(s) where it is licensed/approved/registered?
                </Form.Label>
                <Col>
                    <Field.Switch
                        type="switch"
                        label="Must comply with the Travel Rule"
                        register={register}
                        name="must_comply_travel_rule"
                    />
                </Col>
            </Form.Group>
            <Regulations register={register} name="applicable_regulations" control={control} />
            <Form.Group as={Row} className="mb-3">
                <Form.Label column sm="12" className="fw-normal">
                    What is the minimum threshold for Travel Rule compliance?
                </Form.Label>
                <Col sm="7">
                    <Field.Input type="number" register={register} name="compliance_threshold" />
                </Col>
                <Col sm="3">
                    <Field.Select register={register} name="compliance_threshold_currency">
                        <Currencies />
                    </Field.Select>
                </Col>
                <Form.Text>
                    The minimum threshold above which your organization is required to collect/send Travel Rule
                    information.
                </Form.Text>
            </Form.Group>

            <h4>Data Protection Policies</h4>
            <Form.Group as={Row} className="mb-3" controlId="must_safeguard_pii">
                <Form.Label column sm="12" className="fw-normal">
                    Is your organization required by law to safeguard PII?
                </Form.Label>
                <Col sm="7">
                    <Field.Switch
                        type="switch"
                        label="Must Safeguard PII"
                        register={register}
                        name="must_safeguard_pii"
                    />
                </Col>
            </Form.Group>

            <Form.Group as={Row} className="mb-3" controlId="safeguards_pii">
                <Form.Label column sm="12" className="fw-normal">
                    Does your organization secure and protect PII, including PII received from other VASPs under the
                    Travel Rule?
                </Form.Label>
                <Col sm="7">
                    <Field.Switch type="switch" label="Safeguards PII" register={register} name="safeguards_pii" />
                </Col>
            </Form.Group>
            <div className="mt-3 text-end">
                <ModalCloseButton>
                    <Button variant="danger" className="me-2">
                        Cancel
                    </Button>
                </ModalCloseButton>
                <Button type="submit" disabled={isLoading || !isDirty}>
                    Save
                </Button>
            </div>
        </Form>
    );
}

export default TrixoEditForm;
