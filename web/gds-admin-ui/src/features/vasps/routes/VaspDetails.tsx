import React, { Suspense } from 'react';
import { Col, Row } from 'react-bootstrap';
import { useParams } from 'react-router-dom';

import PageTitle from '@/components/PageTitle';

import AuditLog from '../components/details/AuditLog';
import BasicDetails from '../components/details/BasicDetails/BasicDetails';
import CertificateDetails from '../components/details/CertificateDetails/CertificateDetails';
import Contact from '../components/details/contact';
import EmailLog from '../components/details/EmailLog';
import TrixoQuestionnaire from '../components/details/TrixoQuestionnaire';
import { useGetVasp } from '../services';
import OvalLoader from '@/components/OvalLoader';

const ReviewNotes = React.lazy(() => import('../components/details/ReviewNotes'));

const VaspDetails = () => {
    const params = useParams<{ id: string }>();
    const { data: vasp, isLoading } = useGetVasp({ vaspId: params.id });

    return (
        <>
            <PageTitle
                breadCrumbItems={[
                    { label: 'List', path: '/vasps' },
                    { label: 'Details', path: `/vasps/${params?.id}`, active: true },
                ]}
                title="Registration Details"
            />
            {isLoading ? <OvalLoader /> : null}
            {vasp && (
                <Row>
                    <Col md={6} xl={8} xxl={8}>
                        <BasicDetails data={vasp} />
                        <TrixoQuestionnaire data={vasp?.vasp?.trixo} />
                        <Suspense fallback={<OvalLoader />}>
                            <ReviewNotes />
                        </Suspense>
                    </Col>
                    <Col md={6} xl={4} xxl={4}>
                        <Contact data={vasp?.vasp?.contacts} verifiedContact={vasp?.verified_contacts} />
                        <AuditLog data={vasp?.audit_log} />
                        <EmailLog data={vasp?.email_log} />
                        <CertificateDetails data={vasp?.vasp?.identity_certificate} />
                    </Col>
                </Row>
            )}
        </>
    );
};

export default VaspDetails;
