
import React from 'react'
import { LEGAL_PERSON_NAME_IDENTIFIER_TYPE } from 'constants/index'

function LegalPersonNameIdentifierTypeOptions() {
    return (
        <>
            {
                LEGAL_PERSON_NAME_IDENTIFIER_TYPE.map((type, idx) => <option key={type} value={idx}>{type}</option>)
            }
        </>
    )
}

export default LegalPersonNameIdentifierTypeOptions
