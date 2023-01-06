import React, { Suspense } from 'react';
import { Col, Row } from 'react-bootstrap';

// components
import PageTitle from '@/components/PageTitle';
import { Statistics, Status } from '../components/dashboard';
import { useGetSummary } from '../services';
import { lazyImport } from '@/lib/lazy-import';
import OvalLoader from '@/components/OvalLoader';

const { VaspsByCountryChart } = lazyImport(() => import('../components/dashboard'), 'VaspsByCountryChart');
const { TasksChart } = lazyImport(() => import('../components/dashboard'), 'TasksChart');
const { PendingAndRecentActivity } = lazyImport(() => import('../components/dashboard'), 'PendingAndRecentActivity');

const Dashboard = () => {
    const { data: summary } = useGetSummary();

    return (
        <>
            <PageTitle breadCrumbItems={[{ label: 'Summary', path: '/dashboard', active: true }]} title="Dashboard" />
            <Statistics data={summary} />
            <Row>
                <Col lg={4}>
                    <Status statuses={summary?.statuses} />
                </Col>
                <Col sm={12} lg={8} className="d-flex">
                    <Suspense fallback={<OvalLoader />}>
                        <PendingAndRecentActivity />
                    </Suspense>
                </Col>
            </Row>
            <Row>
                <Col>
                    <Suspense fallback={<OvalLoader />}>
                        <TasksChart />
                    </Suspense>
                </Col>
            </Row>
            <Row>
                <Col>
                    <Suspense fallback={<OvalLoader />}>
                        <VaspsByCountryChart />
                    </Suspense>
                </Col>
            </Row>
        </>
    );
};

export default React.memo(Dashboard);
