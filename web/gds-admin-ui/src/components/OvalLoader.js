
import React from 'react'
import Oval from './Oval'

function OvalLoader({ title, ...rest }) {
    return (
        <div className='text-center flex flex-column justify-content-center' {...rest}>
            <div>
                <Oval width="40" height="40" stroke="#6b7280" />
            </div>
            <div className='mt-1'>
                <small className='block'>{title || 'Loading...'}</small>
            </div>
        </div>
    )
}

export default OvalLoader
