
import React from 'react'
import { Dropdown } from 'react-bootstrap'
import PropTypes from 'prop-types'

function Print({ onPrint }) {
    return (
        <Dropdown.Item onClick={onPrint} data-testid="print-btn">
            <i className="mdi mdi-printer me-1"></i>Print
        </Dropdown.Item>
    )
}

Print.propTypes = {
    onPrint: PropTypes.func.isRequired
}

export default Print