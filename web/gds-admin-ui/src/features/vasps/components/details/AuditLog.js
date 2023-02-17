import { Card } from "react-bootstrap";
import SimpleBar from 'simplebar-react';
import dayjs from "dayjs";
import PropTypes from 'prop-types';
import { StatusLabel } from "@/constants";
import TimelineItem from "@/components/TimelineItem";
import Timeline from "@/components/Timeline";

function sortByTimestamp(data = []) {
    return data.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
}

const AuditLog = ({ data }) => {
    return data ? (
        <Card>
            <Card.Body>
                <h4 className="header-title mb-2">Audit log</h4>

                <SimpleBar style={{ maxHeight: '330px', width: '100%' }}>
                    <Timeline>
                        {
                            sortByTimestamp(data).map((log, idx) => (
                                <TimelineItem key={idx}>
                                    <i className="mdi mdi-file-search-outline bg-info-lighten text-info timeline-icon"></i>
                                    <div className="timeline-item-info mb-3">
                                        <h5 className="text-primary mb-1" data-testid="audit-log-desc" href="/">{log.description}</h5>
                                        <small className="d-block" data-testid="audit-log-state">
                                            <span className="fst-italic" data-testid="audit-log-source">{log.source}</span> from{' '}
                                            <span className="fst-italic" data-testid="audit-log-previous-state">{StatusLabel[log.previous_state]}</span> &rarr;{' '}
                                            <span className="fst-italic" data-testid="audit-log-current-state">{StatusLabel[log.current_state]}</span>
                                        </small>
                                        <small className="d-block" data-testid="audit-log-timestamp">{dayjs(log?.timestamp).toLocaleString()}</small>
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

AuditLog.propTypes = {
    data: PropTypes.arrayOf(PropTypes.object)
}

export default AuditLog;