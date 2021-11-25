import React from 'react';
import { Row, Col } from 'react-bootstrap';
import { useDispatch, useSelector } from 'react-redux';

// components
import PageTitle from 'components/PageTitle';
import { fecthRegistrationsReviews, fetchCertificates, fetchPendingVasps, fetchSummary } from 'redux/dashboard/actions';

import Statistics from './Statistics';
import Status from './Status';
import Tasks from './Tasks';
import TasksChart from './TasksChart';


const ProjectDashboardPage = () => {
    const dispatch = useDispatch();
    const { summary, vasps, isVaspsLoading } = useSelector(state => ({
        summary: state.Summary.data,
        vasps: state.Vasps.data,
        isVaspsLoading: state.Vasps.loading
    }))

    React.useEffect(() => {
        dispatch(fetchCertificates());
        dispatch(fetchPendingVasps());
        dispatch(fetchSummary())
        dispatch(fecthRegistrationsReviews())
    }, [dispatch])

    return (
        <React.Fragment>
            <PageTitle
                breadCrumbItems={[
                    { label: 'Summary', path: '/dashboard', active: true }
                ]}
                title={'Dashboard'}
            />

            <Statistics data={summary} />

            <Row>
                <Col lg={4}>
                    <Status statuses={summary?.statuses} />
                </Col>
                <Col sm={12} lg={8} style={{ overflowY: "auto", height: "100%" }}>
                    {!isVaspsLoading ? <Tasks data={vasps} /> : null}
                </Col>
            </Row>
            <Row>
                <Col>
                    <TasksChart />
                </Col>
            </Row>

        </React.Fragment>
    );
};

export default ProjectDashboardPage;
