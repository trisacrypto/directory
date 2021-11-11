import { Card } from "react-bootstrap";
import SimpleBar from 'simplebar-react';
import TimelineItem from "../../../components/TimelineItem";
import Timeline from '../../../components/Timeline'
import { StatusLabel } from "../../../constants";


const AuditLog = ({ data }) => {
    return data ? (
        <Card>
            <Card.Body>
                <h4 className="header-title mb-2">Audit log</h4>

                <SimpleBar style={{ maxHeight: '330px', width: '100%' }}>
                    <Timeline>
                        {
                            data.map((log, idx) => (
                                <TimelineItem key={idx}>
                                    <i className="mdi mdi-file-search-outline bg-info-lighten text-info timeline-icon"></i>
                                    <div className="timeline-item-info">
                                        <h5 data-testid="audit-log-desc" href="/" style={{ color: '#2e74d6' }}>{log.description}</h5>
                                        <small className="d-block" data-testid="audit-log-state">
                                            <span className="fst-italic" data-testid="audit-log-source">{log.source}</span> from{' '}
                                            <span className="fst-italic" data-testid="audit-log-previous-state">{StatusLabel[log.previous_state]}</span> -{' '}
                                            <span className="fst-italic" data-testid="audit-log-current-state">{StatusLabel[log.current_state]}</span>
                                        </small>
                                        <small className="d-block" data-testid="audit-log-timestamp">{new Date(log.timestamp).toLocaleString()}</small>
                                    </div>
                                </TimelineItem>
                            ))
                        }

                    </Timeline>
                </SimpleBar>
            </Card.Body>
        </Card>
    ) : null;
};

export default AuditLog;