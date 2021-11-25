import React from 'react';
import { Bar } from 'react-chartjs-2';
import { Card } from 'react-bootstrap';
import { useSelector } from 'react-redux';
import { Status, StatusLabel } from 'constants/dashboard';


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


const TasksChart = () => {
    const { reviews, isLoading } = useSelector(state => ({
        reviews: state.Reviews.data,
        isLoading: state.Reviews.loading
    }))

    const getWeeks = () => {
        if (reviews && Array.isArray(reviews.weeks)) {
            return reviews.weeks.map(review => review.week)
        }

        return []
    }

    const getData = (key = '') => {
        if (reviews && Array.isArray(reviews.weeks)) {
            return reviews.weeks.map(week => {
                return week.registrations[key];
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
                label: StatusLabel.SUBMITTED,
                backgroundColor: '#727cf5',
                borderColor: '#727cf5',
                data: getData(Status.SUBMITTED),
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                label: StatusLabel.EMAIL_VERIFIED,
                backgroundColor: '#1abc9c',
                borderColor: '#1abc9c',
                data: getData(Status.EMAIL_VERIFIED),
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                label: StatusLabel.PENDING_REVIEW,
                backgroundColor: '#3498db',
                borderColor: '#3498db',
                data: getData(Status.PENDING_REVIEW),
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                label: StatusLabel.REVIEWED,
                backgroundColor: '#9b59b6',
                borderColor: '#9b59b6',
                data: getData(Status.REVIEWED),
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                label: StatusLabel.ISSUING_CERTIFICATE,
                backgroundColor: '#f1c40f',
                borderColor: '#f1c40f',
                data: getData(Status.ISSUING_CERTIFICATE),
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                label: StatusLabel.VERIFIED,
                backgroundColor: '#e74c3c',
                borderColor: '#e74c3c',
                data: getData(Status.VERIFIED),
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                label: StatusLabel.REJECTED,
                backgroundColor: '#3B3B98',
                borderColor: '#3B3B98',
                data: getData(Status.REJECTED),
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                label: StatusLabel.APPEALED,
                backgroundColor: '#182C61',
                borderColor: '#182C61',
                data: getData(Status.APPEALED),
            },
            {
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                label: StatusLabel.ERRORED,
                backgroundColor: '#9AECDB',
                borderColor: '#9AECDB',
                data: getData(Status.ERRORED),
            }
        ],
    };


    return (
        <Card>
            <Card.Body>
                <h4 className="header-title mb-4">REVIEWS TIMELINE</h4>

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
