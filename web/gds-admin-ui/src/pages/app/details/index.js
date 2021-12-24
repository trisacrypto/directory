import React from 'react';
import { Row, Col } from 'react-bootstrap';

import PageTitle from 'components/PageTitle';
import { useParams } from "react-router-dom"

import Contact from './contact';
import BasicDetails from './BasicDetails';
import CertificateDetails from './CertificateDetails';
import { getVasp } from "services/vasps"
import { useHistory } from 'react-router-dom'
import AuditLog from './AuditLog';
import { useDispatch } from 'react-redux';
import useSafeDispatch from 'hooks/useSafeDispatch';
import { fetchReviewNotesApiResponse } from 'redux/review-notes';
import NProgress from 'nprogress'
import TrixoQuestionnaire from './TrixoQuestionnaire';

const ReviewNotes = React.lazy(() => import('./ReviewNotes'))


const VaspDetails = () => {
    const [vasp, setVasp] = React.useState({});
    const params = useParams();
    const history = useHistory();
    const dispatch = useDispatch()
    const safeDispatch = useSafeDispatch(dispatch)

    React.useEffect(() => {
        NProgress.start()
        if (params && params.id) {
            safeDispatch(fetchReviewNotesApiResponse(params.id))

            getVasp(params.id).then(response => {
                setVasp(response.data)
                NProgress.done()
            }).catch(error => {
                history.push('/not-found', { error: "Could not retrieve VASP record by ID" })
                console.error("[BasicDetails] getVasp", error.message)
                NProgress.done()
            })

        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [params.id, safeDispatch])



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
                        <AuditLog data={vasp?.audit_log || []} />
                        <CertificateDetails data={vasp?.vasp?.identity_certificate} />
                    </Col>
                </Row>
            )}
        </React.Fragment>
    );
};

export default VaspDetails;
