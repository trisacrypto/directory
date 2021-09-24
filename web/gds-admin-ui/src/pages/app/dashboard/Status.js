import React from 'react';
import { Doughnut } from 'react-chartjs-2';
import { Card, Row, Col } from 'react-bootstrap';
import { capitalizeFirstLetter, convertCountsToPercentages } from '../../../utils';

const Status = ({ statuses }) => {
    const colors = ['#0acf97', '#727cf5', '#fa5c7c'];

    const getDonutChartData = () => {
        if (statuses && typeof statuses === "object") {
            return Object.values(convertCountsToPercentages(statuses))
        }

        return []
    }

    const getDonutChartLabels = () => {
        if (statuses && typeof statuses === "object") {
            return Object.keys(statuses).map(status => {
                const status_ = capitalizeFirstLetter(status).split('_')
                return status_.join(' ')
            })
        }

        return []
    }

    const donutChartData = {
        labels: getDonutChartLabels(),
        datasets: [
            {
                data: getDonutChartData(),
                backgroundColor: colors,
                borderColor: 'transparent',
                borderWidth: '3',
            },
        ],
    };

    const donutChartOpts = {
        maintainAspectRatio: false,
        cutoutPercentage: 80,
        legend: {
            display: false,
        },
    };

    return (
        <Card>
            <Card.Body>
                <h4 className="header-title mb-4">Review Speed</h4>

                <div className="my-4 chartjs-chart" style={{ height: '202px' }}>
                    <Doughnut data={donutChartData} options={donutChartOpts} />
                </div>

                <Row className="text-center mt-2 py-2">
                    <Col lg={4}>
                        <i className="mdi mdi-progress-question text-warning mt-3 h3"></i>
                        <h3 className="fw-normal">
                            <span>{statuses?.PENDING_REVIEW + '%'}</span>
                        </h3>
                        <p className="text-muted mb-0">Pending</p>
                    </Col>

                    <Col lg={4}>
                        <i className="mdi mdi-alert-octagram text-danger mt-3 h3"></i>
                        <h3 className="fw-normal">
                            <span>{statuses?.REJECTED + '%'}</span>
                        </h3>
                        <p className="text-muted mb-0">Rejected</p>
                    </Col>

                    <Col lg={4}>
                        <i className="mdi mdi-shield-check text-primary mt-3 h3"></i>
                        <h3 className="fw-normal">
                            <span>{statuses?.VERIFIED + '%'}</span>
                        </h3>
                        <p className="text-muted mb-0"> Verified</p>
                    </Col>
                </Row>
            </Card.Body>
        </Card>
    );
};

export default Status;
