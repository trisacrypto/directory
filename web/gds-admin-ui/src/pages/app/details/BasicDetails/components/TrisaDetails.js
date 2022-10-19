
import React from 'react';
import { formatDisplayedData } from 'utils';
import Copy from 'components/Copy'


function TrisaDetails({ handleTrisaJsonExportClick, data }) {
    return <>
        <h4 className='text-dark mb-0'>TRISA details <button onClick={handleTrisaJsonExportClick} className='mdi mdi-arrow-down-bold-circle-outline border-0 bg-transparent' title="Download as JSON"></button></h4>
        <hr className='my-1' />
        <p className="mb-2 fw-bold">ID: <span className="fw-normal">{formatDisplayedData(data?.vasp?.id)}</span>
            <Copy data={data?.vasp?.id} />
        </p>
        <p className="mb-2 fw-bold">Common name: <span className="fw-normal">{formatDisplayedData(data?.vasp?.common_name)}</span> <Copy data={data?.vasp?.common_name} /> </p>
        <p className="mb-2 fw-bold">Endpoint: <span className="fw-normal">{formatDisplayedData(data?.vasp?.trisa_endpoint)}</span></p>
        <p className="mb-2 fw-bold">Registered directory: <span className="fw-normal">{formatDisplayedData(data?.vasp?.registered_directory)}</span></p>
    </>;
}

export default TrisaDetails;
