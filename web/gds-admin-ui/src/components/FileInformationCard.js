import React from 'react'
import { Card, Col, Row } from 'react-bootstrap'
import PropTypes from 'prop-types'
import { formatBytes, getBase64Size } from 'utils';

function FileInformationCard({ name, file, ext }) {
    const fileSize = React.useCallback(() => file ? formatBytes(getBase64Size(file)) : '', [file])

    return (
        <Card className="mb-1 shadow-none border">
            <div className="p-2">
                <Row className="align-items-center">
                    <Col className="col-auto">
                        <div className="avatar-sm text-break text-center">
                            <span className="avatar-title bg-primary-lighten text-primary rounded p-1">
                                {ext}
                            </span>
                        </div>
                    </Col>
                    <Col className="col ps-0">
                        <a href="/" className="text-muted font-weight-bold">
                            {name}
                        </a>
                        <p className="mb-0">{fileSize()}</p>
                    </Col>
                    <Col className="col-auto" disabled={true}>
                        <a
                            href="/"
                            className="btn btn-link btn-lg text-muted" disabled>
                            <i className="dripicons-download"></i>
                        </a>
                    </Col>
                </Row>
            </div>
        </Card>
    )
}

FileInformationCard.propTypes = {
    file: PropTypes.string
}

export default FileInformationCard
