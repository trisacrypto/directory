import React from "react";
import { Card } from "react-bootstrap";
import SimpleBar from 'simplebar-react';
import TimelineItem from "components/TimelineItem";
import Timeline from 'components/Timeline'
import dayjs from "dayjs";
import PropTypes from 'prop-types';

function sortByTimestamp(data = []) {
    return data.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
}

const EmailLog = ({ data }) => {
    return data ? (
        <Card>
            <Card.Body>
                <h4 className="header-title mb-2">Email log</h4>

                {
                    Array.isArray(data) && data.length ? (

                        <SimpleBar style={{ maxHeight: '330px', width: '100%' }}>
                            <Timeline>
                                {
                                    sortByTimestamp(data).map((log, idx) => (
                                        <TimelineItem key={idx}>
                                            <i className="mdi mdi-file-search-outline bg-info-lighten text-info timeline-icon"></i>
                                            <div className="timeline-item-info mb-3">
                                                <h5 className="text-primary mb-1" data-testid="email-log-reason" href="/">{log.reason}</h5>
                                                <small className="d-block" data-testid="email-log-contact-line">
                                                    contact: <span className="fst-italic" data-testid="email-log-contact">{log.contact}</span>
                                                </small>
                                                <small data-testid="email-log-subject-line">
                                                    subject: <span className="fst-italic" data-testid="email-log-subject">{log.subject}</span>
                                                </small>
                                                <small className="d-block" data-testid="email-log-timestamp">{dayjs(log?.timestamp).toLocaleString()}</small>
                                            </div>
                                        </TimelineItem>
                                    ))
                                }

                            </Timeline>
                        </SimpleBar>
                    ) : (
                        <p className="fst-italic fs-6" data-testid="no-data">No Email Log</p>
                    )
                }
            </Card.Body>
        </Card>
    ) : null;
};

EmailLog.propTypes = {
    data: PropTypes.arrayOf(PropTypes.object)
}

export default EmailLog;