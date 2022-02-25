
import React from 'react'
import { isoCountries } from 'utils/country'

function CountryOptions() {
    return (
        <>
            {
                Object.entries(isoCountries).map(([k, v]) => <option key={k} value={k}>{v}</option>)
            }
        </>
    )
}

export default CountryOptions
