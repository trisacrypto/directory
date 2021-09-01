// @flow
import React from 'react';
import { Row, Col } from 'react-bootstrap';
import { useDispatch } from 'react-redux';

// components
import PageTitle from '../../../components/PageTitle';
import { fetchCertificates, fetchSummary, fetchVasps } from '../../../redux/dashboard/actions';

import Statistics from './Statistics';
import Status from './Status';
import Tasks from './Tasks';


const ProjectDashboardPage = (): React$Element<React$FragmentType> => {
    const dispatch = useDispatch();

    React.useEffect(() => {
        dispatch(fetchCertificates());
        dispatch(fetchVasps());
        dispatch(fetchSummary())
    }, [dispatch])

    return (
        <React.Fragment>
            <PageTitle
                breadCrumbItems={[]}
                title={'Dashboard'}
            />

            <Statistics />

            <Row style={{ height: "500px" }}>
                <Col lg={4}>
                    <Status />
                </Col>
                <Col lg={8} style={{ overflowY: "scroll", height: "100%" }}>
                    <Tasks />
                </Col>
            </Row>
        </React.Fragment>
    );
};

export default ProjectDashboardPage;
