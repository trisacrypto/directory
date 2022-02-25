
import React from 'react'
import { AddressTypeHeaders } from 'constants/index'

function AddressTypeOptions() {
    return (
        <>
            {
                Object.entries(AddressTypeHeaders).map(([k, v]) => <option key={k} value={k}>{v}</option>)
            }
        </>
    )
}

export default AddressTypeOptions
