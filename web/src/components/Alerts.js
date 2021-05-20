import React from 'react';
import Alert from 'react-bootstrap/Alert';

const Alerts = ({ alerts, onDismiss }) => {
  const renderedAlerts = alerts.map((alert, idx) => (
    <Alert key={idx} onClose={() => onDismiss(idx)} variant={alert.variant} dismissible>
      {alert.message}
    </Alert>
  ));
  return <div id="alerts">{renderedAlerts}</div>
}

export default Alerts;