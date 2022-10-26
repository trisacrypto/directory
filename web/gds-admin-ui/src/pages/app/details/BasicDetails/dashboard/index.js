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
import { getPendingVaspsData, getPendingVaspsError, getPendingVaspsLoadingState, getSummaryData } from 'redux/selectors/dashboard';
import useSafeDispatch from 'hooks/useSafeDispatch';


const ProjectDashboardPage = () => {
    const _dispatch = useDispatch();
    const safeDispatch = useSafeDispatch(_dispatch)
    const summary = useSelector(getSummaryData)
    const vasps = useSelector(getPendingVaspsData)
    const isVaspsLoading = useSelector(getPendingVaspsLoadingState)
    const pendingVaspsError = useSelector(getPendingVaspsError)

    React.useEffect(() => {
        safeDispatch(fetchCertificates());
        safeDispatch(fetchPendingVasps());
        safeDispatch(fetchSummary())
        safeDispatch(fecthRegistrationsReviews())
    }, [safeDispatch])

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
                <Col sm={12} lg={8} className='d-flex'>
                    <Tasks data={vasps} isLoading={isVaspsLoading} error={pendingVaspsError} />
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

export default React.memo(ProjectDashboardPage);
