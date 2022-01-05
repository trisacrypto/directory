
import { BUSINESS_CATEGORY } from 'constants/basic-details'
import React from 'react'

function BusinessCategory() {
    return (
        <>
            <option value={''}></option>
            {
                Object.entries(BUSINESS_CATEGORY).map(([k, v]) => (
                    <option value={k} key={k}>{v}</option>
                ))
            }
        </>
    )
}

export default BusinessCategory
