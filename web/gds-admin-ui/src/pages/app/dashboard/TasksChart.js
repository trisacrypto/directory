// @flow
import React from 'react';
import { Bar } from 'react-chartjs-2';
import { Card } from 'react-bootstrap';
import { useSelector } from 'react-redux';
<<<<<<< HEAD
import { Status, StatusLabel } from '../../../constants/dashboard';


const barChartOpts = {
    maintainAspectRatio: false,
    legend: {
        display: true,
    },
    tooltips: {
        intersect: false,
    },
    hover: {
        intersect: true,
    },
    plugins: {
        filler: {
            propagate: false,
        },
    },
    scales: {
        xAxes: [
            {
                reverse: true,
                gridLines: {
                    color: 'rgba(0,0,0,0.05)',
                },
                stacked: true
            },
        ],
        yAxes: [
            {
                ticks: {
                    stepSize: 10,
                    display: false,
                },
                min: 10,
                max: 100,
                display: true,
                borderDash: [5, 5],
                gridLines: {
                    color: 'rgba(0,0,0,0)',
                    fontColor: '#fff',
                },
                stacked: true
            },
        ],
    },
};
=======
import { STACKED_BAR_LABEL } from '../../../constants/dashboard';
>>>>>>> feat: add reviews timeline graph


const TasksChart = (): React$Element<any> => {
    const { reviews, isLoading } = useSelector(state => ({
        reviews: state.Reviews.data,
        isLoading: state.Reviews.loading
    }))

    const getWeeks = () => {
        if (reviews && Array.isArray(reviews)) {
            return reviews.map(review => review.week)
        }

        return []
    }

    const getData = (key = '') => {
        if (reviews && Array.isArray(reviews)) {
            return reviews.map(review => {
                return review.registrations[key];
            })
        }

        return []
    }

    const barChartData = {
        labels: [
            ...getWeeks()
        ],
        datasets: [
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.SUBMITTED,
                backgroundColor: '#727cf5',
                borderColor: '#727cf5',
                data: getData(Status.SUBMITTED),
=======
                label: STACKED_BAR_LABEL.SUBMITTED,
                backgroundColor: '#727cf5',
                borderColor: '#727cf5',
                data: getData(STACKED_BAR_LABEL.SUBMITTED),
>>>>>>> feat: add reviews timeline graph
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.EMAIL_VERIFIED,
                backgroundColor: '#1abc9c',
                borderColor: '#1abc9c',
                data: getData(Status.EMAIL_VERIFIED),
=======
                label: STACKED_BAR_LABEL.EMAIL_VERIFIED,
                backgroundColor: '#1abc9c',
                borderColor: '#1abc9c',
                data: getData(STACKED_BAR_LABEL.EMAIL_VERIFIED),
>>>>>>> feat: add reviews timeline graph
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.PENDING_REVIEW,
                backgroundColor: '#3498db',
                borderColor: '#3498db',
                data: getData(Status.PENDING_REVIEW),
=======
                label: STACKED_BAR_LABEL.PENDING_REVIEW,
                backgroundColor: '#3498db',
                borderColor: '#3498db',
                data: getData(STACKED_BAR_LABEL.PENDING_REVIEW),
>>>>>>> feat: add reviews timeline graph
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.REVIEWED,
                backgroundColor: '#9b59b6',
                borderColor: '#9b59b6',
                data: getData(Status.REVIEWED),
=======
                label: STACKED_BAR_LABEL.REVIEWED,
                backgroundColor: '#9b59b6',
                borderColor: '#9b59b6',
                data: getData(STACKED_BAR_LABEL.REVIEWED),
>>>>>>> feat: add reviews timeline graph
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.ISSUING_CERTIFICATE,
                backgroundColor: '#f1c40f',
                borderColor: '#f1c40f',
                data: getData(Status.ISSUING_CERTIFICATE),
=======
                label: STACKED_BAR_LABEL.ISSUING_CERTIFICATE,
                backgroundColor: '#f1c40f',
                borderColor: '#f1c40f',
                data: getData(STACKED_BAR_LABEL.ISSUING_CERTIFICATE),
>>>>>>> feat: add reviews timeline graph
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.VERIFIED,
                backgroundColor: '#e74c3c',
                borderColor: '#e74c3c',
                data: getData(Status.VERIFIED),
=======
                label: STACKED_BAR_LABEL.VERIFIED,
                backgroundColor: '#e74c3c',
                borderColor: '#e74c3c',
                data: getData(STACKED_BAR_LABEL.VERIFIED),
>>>>>>> feat: add reviews timeline graph
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.REJECTED,
                backgroundColor: '#3B3B98',
                borderColor: '#3B3B98',
                data: getData(Status.REJECTED),
=======
                label: STACKED_BAR_LABEL.REJECTED,
                backgroundColor: '#3B3B98',
                borderColor: '#3B3B98',
                data: getData(STACKED_BAR_LABEL.REJECTED),
>>>>>>> feat: add reviews timeline graph
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.APPEALED,
                backgroundColor: '#182C61',
                borderColor: '#182C61',
                data: getData(Status.APPEALED),
=======
                label: STACKED_BAR_LABEL.APPEALED,
                backgroundColor: '#182C61',
                borderColor: '#182C61',
                data: getData(STACKED_BAR_LABEL.APPEALED),
>>>>>>> feat: add reviews timeline graph
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
<<<<<<< HEAD
                label: StatusLabel.ERRORED,
                backgroundColor: '#9AECDB',
                borderColor: '#9AECDB',
                data: getData(Status.ERRORED),
=======
                label: STACKED_BAR_LABEL.ERRORED,
                backgroundColor: '#9AECDB',
                borderColor: '#9AECDB',
                data: getData(STACKED_BAR_LABEL.ERRORED),
>>>>>>> feat: add reviews timeline graph
            }
        ],
    };

<<<<<<< HEAD
=======
    const barChartOpts = {
        maintainAspectRatio: false,
        legend: {
            display: true,
        },
        tooltips: {
            intersect: false,
        },
        hover: {
            intersect: true,
        },
        plugins: {
            filler: {
                propagate: false,
            },
        },
        scales: {
            xAxes: [
                {
                    reverse: true,
                    gridLines: {
                        color: 'rgba(0,0,0,0.05)',
                    },
                    stacked: true
                },
            ],
            yAxes: [
                {
                    ticks: {
                        stepSize: 10,
                        display: false,
                    },
                    min: 10,
                    max: 100,
                    display: true,
                    borderDash: [5, 5],
                    gridLines: {
                        color: 'rgba(0,0,0,0)',
                        fontColor: '#fff',
                    },
                    stacked: true
                },
            ],
        },
    };
>>>>>>> feat: add reviews timeline graph

    return (
        <Card>
            <Card.Body>
<<<<<<< HEAD
                <h4 className="header-title mb-4">REVIEWS TIMELINE</h4>
=======
                <h4 className="header-title mb-4">REGISTRATION</h4>
>>>>>>> feat: add reviews timeline graph

                <div dir="ltr">
                    <div style={{ height: '320px' }} className="mt-3 chartjs-chart">
                        {
                            !isLoading && <Bar data={barChartData} options={barChartOpts} />
                        }
                    </div>
                </div>
                <small>
                    ⋆ click on the elements of the legend to filter accordingly
                </small>
            </Card.Body >
        </Card >
    );
};

export default TasksChart;
