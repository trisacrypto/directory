// @flow
import React from 'react';
import { Row, Col } from 'react-bootstrap';
import { useDispatch, useSelector } from 'react-redux';
import config from '../../../config'

// components
import PageTitle from '../../../components/PageTitle';
import { fecthRegistrationsReviews, fetchCertificates, fetchPendingVasps, fetchSummary } from '../../../redux/dashboard/actions';

import Statistics from './Statistics';
import Status from './Status';
import Tasks from './Tasks';
import TasksChart from './TasksChart';
import { ENVIRONMENT } from '../../../constants';


const ProjectDashboardPage = (): React$Element<React$FragmentType> => {
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

            <Row style={{ height: "500px" }}>
                <Col lg={4}>
                    <Status statuses={summary.statuses} />
                </Col>
                <Col lg={8} style={{ overflowY: "scroll", height: "100%" }}>
                    {!isVaspsLoading ? <Tasks data={vasps} /> : null}
                </Col>
            </Row>
            {
                config.ENVIRONMENT === ENVIRONMENT.PROD ? null : (
                    <Row>
                        <Col>
                            <TasksChart />
                        </Col>
                    </Row>)
            }
        </React.Fragment>
    );
};

export default ProjectDashboardPage;
