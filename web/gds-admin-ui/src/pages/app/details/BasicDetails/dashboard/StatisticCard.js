import PropTypes from 'prop-types';
import { Card } from 'react-bootstrap';

function StatisticCard({ count, title, icon }) {
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

StatisticCard.propTypes = {
  count: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  icon: PropTypes.node.isRequired,
};

export default StatisticCard;
