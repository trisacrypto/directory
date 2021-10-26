
import React from 'react'
import { Col } from 'react-bootstrap'
import { NAME_IDENTIFIER_TYPE } from '../../../../../constants/basic-details'
import { formatDisplayedData } from '../../../../../utils'

function Name({ data }) {
    const names = React.useMemo(() => data && [].concat(...Object.values(data)), [data])

    return (
        <Col className="mt-3">
            <p className="mb-1 fw-bold">Name identifiers</p>
            <hr className='mb-1 mt-0' />
            <ul>
                {
                    names?.map((value, index) => value ?
                        (
                            <li key={index}> <span className='badge bg-primary rounded-pill'>{NAME_IDENTIFIER_TYPE[value?.legal_person_name_identifier_type]}</span> {formatDisplayedData(value?.legal_person_name)} </li>
                        ) : null
                    )
                }
            </ul>
        </Col>
    )
}

export default Name
