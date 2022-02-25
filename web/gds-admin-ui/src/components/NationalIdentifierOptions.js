
import React from 'react'

function NationalIdentifierOptions() {
    return (
        <>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_ARNU">Alien Registration Number</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_CCPT">Passport Number</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_RAID">Registration Authority Identifier</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_DRLC">Driver's License Number</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_FIIN">Foreign Investment Identity Number</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_TXID">Tax Identification Number</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_SOCS">Social Security Number</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_IDCD">Identity Card Number</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_LEIX">Legal Entity Identifier (LEI)</option>
            <option value="NATIONAL_IDENTIFIER_TYPE_CODE_MISC">Unspecified</option>
        </>
    )
}

export default NationalIdentifierOptions