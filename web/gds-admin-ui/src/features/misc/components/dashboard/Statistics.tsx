import { Card, Col, Row } from 'react-bootstrap';

import StatisticCard from './StatisticCard';
import { Summary } from '../../types';

type StatisticCardProps = {
    data: Summary;
};

const Statistics = ({ data }: StatisticCardProps) => {
    return (
        <Row>
            <Col>
                <Card className="widget-inline">
                    <Card.Body className="p-0">
                        <Row className="g-0">
                            <Col sm={6} xl={3}>
                                <StatisticCard
                                    title="All VASPs"
                                    count={data?.vasps_count}
                                    icon={<i className="dripicons-briefcase text-muted font-24" />}
                                />
                            </Col>

                            <Col sm={6} xl={3} className="border-start">
                                <StatisticCard
                                    title="Pending Registrations"
                                    count={data?.pending_registrations}
                                    icon={<i className="dripicons-checklist text-muted font-24" />}
                                />
                            </Col>

                            <Col sm={6} xl={3} className="border-start">
                                <StatisticCard
                                    title="Verified Contacts"
                                    count={data?.verified_contacts}
                                    icon={<i className="dripicons-user-group text-muted font-24" />}
                                />
                            </Col>

                            <Col sm={6} xl={3} className="border-start">
                                <StatisticCard
                                    title="Certificates Issued"
                                    count={data?.certificates_issued}
                                    icon={<i className="dripicons-copy text-muted font-24" />}
                                />
                            </Col>
                        </Row>
                    </Card.Body>
                </Card>
            </Col>
        </Row>
    );
};

export default Statistics;
