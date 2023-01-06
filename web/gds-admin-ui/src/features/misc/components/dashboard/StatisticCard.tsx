import { ReactNode } from 'react';
import { Card } from 'react-bootstrap';

type StatisticCardProps = {
    count: number;
    title: string;
    icon: ReactNode;
};

function StatisticCard({ count, title, icon }: StatisticCardProps) {
    return (
        <Card className="shadow-none m-0">
            <Card.Body className="text-center">
                {icon}
                <h3>{count || 0}</h3>
                <p className="text-muted font-15 mb-0">{title}</p>
            </Card.Body>
        </Card>
    );
}

export default StatisticCard;
