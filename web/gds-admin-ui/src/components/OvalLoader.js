
import React from 'react'
import Oval from './Oval'

function OvalLoader() {
    return (
        <div className='text-center flex flex-column justify-content-center'>
            <div>
                <Oval width="40" height="40" stroke="#6b7280" />
            </div>
            <div className='mt-1'>
                <small className='block'>Loading...</small>
            </div>
        </div>
    )
}

export default OvalLoader
