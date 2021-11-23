import React from 'react';
import { Doughnut } from 'react-chartjs-2';
import { Card, Row, Col } from 'react-bootstrap';
import { capitalizeFirstLetter, getRatios } from 'utils';
import { Status as STATUS } from 'constants/index';
import Loader from 'components/Loader'
import PropTypes from 'prop-types';

const Status = ({ statuses }) => {
    const colors = ['#0d6efd', '#dc3545', '#ffc107'];

    const getStatusesCounts = React.useCallback(() => {
        const initialValue = { 'VERIFIED': 0, 'REJECTED': 0, 'PENDING_REVIEW': 0 }
        const reducer = (counts, status) => {
            switch (status[0]) {
                case STATUS.VERIFIED:
                    counts[STATUS.VERIFIED] += status[1]
                    break;
                case STATUS.ERRORED:
                case STATUS.REJECTED:
                    counts[STATUS.REJECTED] += status[1]
                    break;
                default:
                    counts[STATUS.PENDING_REVIEW] += status[1]
                    break;
            }
            return counts
        }
        try {
            return Object.entries(statuses).reduce(reducer, initialValue);
        } catch (error) {
            throw error
        }
    }, [statuses])

    const statusRatios = () => {
        if (statuses && typeof statuses === "object") {
            const statusesCounts = getStatusesCounts(statuses)
            return getRatios(statusesCounts)
        }

        return {}
    }

    const getDonutChartData = () => statuses ? Object.values(getStatusesCounts()) : []
    const statusPercents = () => Object.fromEntries(Object.entries(statusRatios()).map(([key, val]) => [key, val * 100.0]))

    const getDonutChartLabels = () => {
        if (statuses && typeof statuses === "object") {
            return Object.keys(getStatusesCounts()).map(status => {
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
                <h4 className="header-title mb-4">Registration Statuses</h4>

                {
                    !statuses ? <Loader width={50} /> : <>
                        <div className="my-4 chartjs-chart" style={{ height: '202px' }}>
                            <Doughnut data={donutChartData} options={donutChartOpts} />
                        </div>

                        <Row className="text-center mt-2 py-2">
                            <Col lg={4}>
                                <i className="mdi mdi-progress-question text-warning mt-3 h4"></i>
                                <h4 className="fw-normal">
                                    <span>{statusPercents()[STATUS.PENDING_REVIEW] + '%'}</span>
                                </h4>
                                <p className="text-muted mb-0 fs-6">Pending</p>
                            </Col>

                            <Col lg={4}>
                                <i className="mdi mdi-alert-octagram text-danger mt-3 h4"></i>
                                <h4 className="fw-normal">
                                    <span>{statusPercents()[STATUS.REJECTED] + '%'}</span>
                                </h4>
                                <p className="text-muted mb-0 fs-6">Rejected</p>
                            </Col>
                            <Col lg={4}>
                                <i className="mdi mdi-shield-check text-primary mt-3 h4"></i>
                                <h4 className="fw-normal">
                                    <span>{statusPercents()[STATUS.VERIFIED] + '%'}</span>
                                </h4>
                                <p className="text-muted mb-0 fs-6"> Verified</p>
                            </Col>

                        </Row>
                    </>
                }

            </Card.Body>
        </Card>
    );
};

Status.propTypes = {
    statuses: PropTypes.object
}

export default Status;
