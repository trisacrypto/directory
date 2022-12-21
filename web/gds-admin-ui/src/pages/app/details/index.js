import React from 'react';
import { Row, Col } from 'react-bootstrap';

import PageTitle from 'components/PageTitle';
import { useHistory, useParams } from "react-router-dom"

import Contact from './contact';
import BasicDetails from './BasicDetails';
import CertificateDetails from './CertificateDetails';
import AuditLog from './AuditLog';
import { useDispatch } from 'react-redux';
import useSafeDispatch from 'hooks/useSafeDispatch';
import TrixoQuestionnaire from './TrixoQuestionnaire';
import { useSelector } from 'react-redux';
import { getVaspDetails } from 'redux/selectors';
import { fetchVaspDetailsApiResponse } from 'redux/vasp-details';
import EmailLog from './EmailLog';


const ReviewNotes = React.lazy(() => import('./ReviewNotes'))

const VaspDetails = () => {
    const params = useParams();
    const vasp = useSelector(getVaspDetails)
    const dispatch = useDispatch()
    const safeDispatch = useSafeDispatch(dispatch)
    const history = useHistory()

    React.useEffect(() => {
        if (params && params.id) {
            safeDispatch(fetchVaspDetailsApiResponse(params.id, history))
        }

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [params.id, safeDispatch])

    // TO-DO: should review later by adding error page
    return (
        <React.Fragment>
            <PageTitle
                breadCrumbItems={[
                    { label: 'List', path: '/vasps' },
                    { label: 'Details', path: `/vasps/${params?.id}`, active: true }
                ]}
                title={'Registration Details'}
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
        </React.Fragment>
    );
};

export default VaspDetails;
