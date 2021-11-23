
import React from 'react'
import { Col, Row } from 'react-bootstrap'
import { AddressTypeHeaders } from 'constants/index'
import PropTypes from 'prop-types';


export const renderLines = (address) => (
    <address data-testid="addressLine">
        {
            address.address_line.map((addressLine, index) => addressLine &&
                <div key={index} >
                    {addressLine}{' '}
                </div>
            )
        }
        <div>{address?.country}</div>
    </address>
)

export const renderField = (address) => (
    <address data-testid="addressField">
        {address.sub_department ? <>{address?.sub_department} <br /></> : null}
        {address.department ? <>{address?.department} <br /></> : null}
        {address.building_number} {address?.street_name}<br />

        {address.post_box ? <>P.O. Box: {address?.post_box} <br /></> : null}

        {address.floor || address.room || address.building_name ? <>{address?.floor} {address?.room} {address?.building_name} <br /></> : null}

        {address.district_name ? <>{address?.district_name} <br /></> : null}

        {address.town_name || address.town_location_name || address.country_sub_division ? <>{address?.town_name} {address?.town_location_name} {address?.country_sub_division} {address?.post_code}  <br /></> : null}
        {address?.country}
    </address>
)



function Geographic({ data }) {

    const isValidIvmsAddress = (address) => {
        if (address) {
            return !!(address.country && address.address_type)
        }
        return false;
    }

    const hasAddressLine = (address) => {
        if (isValidIvmsAddress(address)) {
            return Array.isArray(address.address_line) && address.address_line.length > 0
        }
        return false;
    }

    const hasAddressField = (address) => {
        if (isValidIvmsAddress(address) && !hasAddressLine(address)) {
            return !!(address.street_name && (address.building_number || address.building_name))
        }
        return false
    }

    const hasAddressFieldAndLine = (address) => {
        if (hasAddressField(address) && hasAddressLine(address)) {
            console.warn("[ERROR]", "cannot render address")
            return true
        }
        return false
    }

    const renderAddress = (address) => {
        if (hasAddressFieldAndLine(address)) {
            console.warn("[ERROR]", "invalid address with both fields and lines");
            return <div>Invalid Address</div>
        }

        if (hasAddressLine(address)) {
            return renderLines(address)
        }

        if (hasAddressField(address)) {
            return renderField(address)
        }

        console.warn("[ERROR]", "could not render address")
        return <div>Unparseable Address</div>
    }



    return (
        <Col lg={6}>
            {
                data && data.map((address, index) => (
                    <Row key={index}>
                        <p className="mb-1 fw-bold" data-testid="addressType">{AddressTypeHeaders[address?.address_type]} Address:</p>
                        {renderAddress(address)}
                    </Row>
                ))
            }
        </Col>
    )
}

Geographic.propTypes = {
    data: PropTypes.arrayOf(PropTypes.object).isRequired
}

export default Geographic
