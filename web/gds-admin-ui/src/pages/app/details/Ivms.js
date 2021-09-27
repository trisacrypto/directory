
import React from 'react'
import { Card, Col, Row } from 'react-bootstrap';
import { formatDisplayedData } from "../../../utils"

function Ivms({ data }) {

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">IVMS 1010</h4>
                <p className="fw-bold mb-2">Country Of Registration: <span className="fw-normal">{formatDisplayedData(data?.country_of_registration)}</span></p>
                <p className="fw-bold mb-2">Customer Number: <span className="fw-normal">{formatDisplayedData(data?.customer_number)}</span></p>
                <Row>
                    <Col className="mt-3">
                        <p className="fw-bold mb-2">Geographic Adressess</p>
                        <hr />
                        {
                            data && data.geographic_addresses && data.geographic_addresses.map((address, index) => (
                                <Row key={index}>
                                    <Col xl={6}>
                                        <p className="fw-bold mb-2">Adress Line: <span className="fw-normal">{formatDisplayedData(address?.address_line).replace(/,/g, " ")}</span></p>
                                        <p className="fw-bold mb-2">Adress Type: <span className="fw-normal">{formatDisplayedData(address?.address_type)}</span></p>
                                        <p className="fw-bold mb-2">Building Name: <span className="fw-normal">{formatDisplayedData(address?.building_name)}</span></p>
                                        <p className="fw-bold mb-2">Building Number: <span className="fw-normal">{formatDisplayedData(address?.building_number)}</span></p>
                                        <p className="fw-bold mb-2">Country: <span className="fw-normal">{formatDisplayedData(address?.country)}</span></p>
                                        <p className="fw-bold mb-2">Country Sub-Division: <span className="fw-normal">{formatDisplayedData(address?.country_sub_division)}</span></p>
                                        <p className="fw-bold mb-2">Department: <span className="fw-normal">{formatDisplayedData(address?.department)}</span></p>
                                        <p className="fw-bold mb-2">Sub Department: <span className="fw-normal">{formatDisplayedData(address?.sub_department)}</span></p>
                                    </Col>
                                    <Col xl={6}>
                                        <p className="fw-bold mb-2">District Name: <span className="fw-normal">{formatDisplayedData(address?.district_name)}</span></p>
                                        <p className="fw-bold mb-2">Floor: <span className="fw-normal">{formatDisplayedData(address?.floor)}</span></p>
                                        <p className="fw-bold mb-2">Post Box: <span className="fw-normal">{formatDisplayedData(address?.post_box)}</span></p>
                                        <p className="fw-bold mb-2">Post Code: <span className="fw-normal">{formatDisplayedData(address?.post_code)}</span></p>
                                        <p className="fw-bold mb-2">Room: <span className="fw-normal">{formatDisplayedData(address?.room)}</span></p>
                                        <p className="fw-bold mb-2">Street Name: <span className="fw-normal">{formatDisplayedData(address?.street_name)}</span></p>
                                        <p className="fw-bold mb-2">Town Location Name: <span className="fw-normal">{formatDisplayedData(address?.town_location_name)}</span></p>
                                        <p className="fw-bold mb-2">Town Name: <span className="fw-normal">{formatDisplayedData(address?.town_name)}</span></p>

                                    </Col>
                                </Row>
                            ))
                        }
                    </Col>
                    <Col className="mt-3">
                        <p className="fw-bold mb-2">Name</p>
                        <hr />
                        <Row>
                            {
                                data && data.name && data.name.local_name_identifiers.length ? (
                                    <Col xl={6}>
                                        <p className="fw-bold mb-2">Local Name Identifiers</p>
                                        <hr />
                                        {
                                            data && data.name && data.name.local_name_identifiers ? data.name.local_name_identifiers.map((identifier, index) => (
                                                <div key={index}>
                                                    <p className="fw-bold mb-2">Legal Person Name: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name)}</span></p>
                                                    <p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name_identifier_type)}</span></p>
                                                </div>
                                            )) : (<p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">"N/A"</span></p>)
                                        }
                                    </Col>
                                ) : (
                                    <Col>
                                        <p className="fw-bold mb-2">Local Name Identifiers: <span className="fw-normal">{formatDisplayedData(data?.name?.local_name_identifiers)}</span></p>
                                    </Col>
                                )
                            }
                            {
                                data && data.name && data.name.phonetic_name_identifiers.length ? (
                                    <Col xl={6}>
                                        <p className="fw-bold mb-2">Phonetic Name Identifiers:</p>
                                        <hr />
                                        {
                                            data && data.name && data.name.phonetic_name_identifiers ? data.name.phonetic_name_identifiers.map((identifier, index) => (
                                                <div key={index}>
                                                    <p className="fw-bold mb-2">Legal Person Name: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name)}</span></p>
                                                    <p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name_identifier_type)}</span></p>
                                                </div>
                                            )) : (<p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">"N/A"</span></p>)
                                        }
                                    </Col>
                                ) : (
                                    <Col>
                                        <p className="fw-bold mb-2">Phonetic Name Identifiers: <span className="fw-normal">{formatDisplayedData(data?.name?.phonetic_name_identifiers)}</span></p>
                                    </Col>
                                )
                            }


                            {
                                data && data.name && data.name.name_identifiers.length ? (
                                    <Col xl={6}>
                                        <p className="fw-bold mb-2">Name Identifiers:</p>
                                        <hr />
                                        {
                                            data && data.name && data.name.name_identifiers ? data.name.name_identifiers.map((identifier, index) => (
                                                <div key={index}>
                                                    <p className="fw-bold mb-2">Legal Person Name: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name)}</span></p>
                                                    <p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name_identifier_type)}</span></p>
                                                </div>
                                            )) : (
                                                <Col>
                                                    <p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">"N/A"</span></p>
                                                </Col>
                                            )
                                        }
                                    </Col>
                                ) : (
                                    <Col>
                                        <p className="fw-bold mb-2">Name Identifiers: <span className="fw-normal">{formatDisplayedData(data?.name?.name_identifiers)}</span></p>
                                    </Col>
                                )
                            }

                        </Row>
                    </Col>
                    <>
                        {
                            data && data.national_identification ? (
                                <Col xl={12} className="mt-2">
                                    <p className="fw-bold mb-2">National Identification</p>
                                    <hr />
                                    <p className="fw-bold mb-2">Country of Issue: <span className="fw-normal">{formatDisplayedData(data?.national_identification?.country_of_issue)}</span></p>
                                    <p className="fw-bold mb-2">National Identifier: <span className="fw-normal">{formatDisplayedData(data?.national_identification?.national_identifier)}</span></p>
                                    <p className="fw-bold mb-2">National Identification Type: <span className="fw-normal">{formatDisplayedData(data?.national_identification?.national_identifier_type)}</span></p>
                                    <p className="fw-bold mb-2">Registration Authority: <span className="fw-normal">{formatDisplayedData(data?.national_identification?.registration_authority)}</span></p>
                                </Col>

                            ) : (
                                <Col>
                                    <p className="fw-bold mb-2">National Identification: <span className="fw-normal">{formatDisplayedData(data?.national_identification)}</span></p>
                                </Col>
                            )
                        }
                    </>
                </Row>

            </Card.Body>
        </Card>
    )
}

export default Ivms
