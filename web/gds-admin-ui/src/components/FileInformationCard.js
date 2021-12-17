import React from 'react'
import { Card, Col, Row } from 'react-bootstrap'
import PropTypes from 'prop-types'
import { formatBytes, getBase64Size } from 'utils';

function FileInformationCard({ name, file, ext, onDownload }) {
    const fileSize = React.useCallback(() => file ? formatBytes(getBase64Size(file)) : '', [file])

    return (
        <Card className="mb-1 shadow-none border">
            <div className="py-1 px-2">
                <Row className="align-items-center">
                    <Col className="col-auto">
                        <div className="avatar-sm text-break text-center">
                            <span className="avatar-title bg-primary-lighten text-primary rounded p-1">
                                {ext}
                            </span>
                        </div>
                    </Col>
                    <Col className="col ps-0 cursor-pointer">
                        <h6 className="text-muted font-weight-bold m-0">
                            {name}
                        </h6>
                        <p className="mb-0">{fileSize()}</p>
                    </Col>
                    <Col onClick={onDownload} className="col-auto" disabled={!file} title='download'>
                        <p
                            className="btn btn-link btn-lg text-muted m-0 p-0" disabled>
                            <i className="dripicons-download"></i>
                        </p>
                    </Col>
                </Row>
            </div>
        </Card>
    )
}

FileInformationCard.propTypes = {
    file: PropTypes.string,
    name: PropTypes.string,
    ext: PropTypes.oneOf(['.PEM', '.GZ'])
}

export default FileInformationCard
