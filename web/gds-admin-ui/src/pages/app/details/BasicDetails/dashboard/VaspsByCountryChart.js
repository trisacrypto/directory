import React from 'react';
import { Bar } from 'react-chartjs-2';
import { Card } from 'react-bootstrap';
import OvalLoader from 'components/OvalLoader';
import { useQuery } from '@tanstack/react-query'
import { APICore } from 'helpers/api/apiCore';
import { isoCountries } from 'utils/country'

const api = new APICore()

const barChartOpts = {
    maintainAspectRatio: false,
    legend: {
        display: false,
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
            },
        ],
        yAxes: [
            {
                ticks: {
                    stepSize: 10,
                    display: true,
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

const getVaspsByCountry = async () => {
    const response = await api.get('/countries')
    return response.data
}

const useGetVaspsByCountry = () => useQuery(['get-vasps-by-country'], getVaspsByCountry)


const VaspsByCountryChart = () => {
    const { data: countries, isLoading } = useGetVaspsByCountry()

    const getIsoCodeEmojies = () => {
        if (countries && Array.isArray(countries)) {
            return countries.map(country => isoCountries[country?.iso_code])
        }
        return []
    }

    const getRegistrations = () => {
        if (countries && Array.isArray(countries)) {
            return countries.map(country => country?.registrations)
        }
        return []
    }

    const barChartData = {
        labels: [
            ...getIsoCodeEmojies()
        ],
        datasets: [
            {
                label: 'Registrations',
                barPercentage: 0.7,
                categoryPercentage: 0.7,
                backgroundColor: '#727cf5',
                borderColor: '#727cf5',
                data: getRegistrations(),
            }
        ],
    };

    return (
        <Card>
            <Card.Body>
                <h4 className="header-title mb-4">Vasps by country</h4>
                {
                    isLoading ? <div><OvalLoader /></div> : (
                        <div dir="ltr">
                            <div style={{ height: '320px' }} className="mt-3 chartjs-chart">
                                <Bar data={barChartData} options={barChartOpts} />
                            </div>
                        </div>
                    )
                }
            </Card.Body >
        </Card >
    );
};

export default VaspsByCountryChart;
