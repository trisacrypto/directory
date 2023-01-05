import React from 'react';
import { Col, Row } from 'react-bootstrap';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useParams } from 'react-router-dom';

import PageTitle from '@/components/PageTitle';
import useSafeDispatch from '@/hooks/useSafeDispatch';
import { getVaspDetails } from '@/redux/selectors';
import { fetchVaspDetailsApiResponse } from '@/redux/vasp-details';

import AuditLog from '../components/details/AuditLog';
import BasicDetails from '../components/details/BasicDetails/BasicDetails';
import CertificateDetails from '../components/details/CertificateDetails/CertificateDetails';
import Contact from '../components/details/contact';
import EmailLog from '../components/details/EmailLog';
import TrixoQuestionnaire from '../components/details/TrixoQuestionnaire';

const ReviewNotes = React.lazy(() => import('../components/details/ReviewNotes'));

const VaspDetails = () => {
    const params = useParams<{ id: string }>();
    const vasp = useSelector(getVaspDetails);
    const dispatch = useDispatch();
    const safeDispatch = useSafeDispatch(dispatch);
    const history = useHistory();

    React.useEffect(() => {
        if (params && params.id) {
            safeDispatch(fetchVaspDetailsApiResponse(params.id, history));
        }

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [params.id, safeDispatch]);

    // TO-DO: should review later by adding error page
    return (
        <>
            <PageTitle
                breadCrumbItems={[
                    { label: 'List', path: '/vasps' },
                    { label: 'Details', path: `/vasps/${params?.id}`, active: true },
                ]}
                title="Registration Details"
            />
            {vasp && (
                <Row>
                    <Col md={6} xl={8} xxl={8}>
                        <BasicDetails data={vasp} />
                        <TrixoQuestionnaire data={vasp?.vasp?.trixo} />
                        <ReviewNotes />
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
