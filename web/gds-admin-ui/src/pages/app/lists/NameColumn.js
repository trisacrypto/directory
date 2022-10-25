import React from 'react';
import { Link } from 'react-router-dom'

const NameColumn = ({ row }) => {
    const id = row?.original?.id || "";

    return (
        <React.Fragment>
            <p className="m-0 d-inline-block align-middle font-16">
                <Link to={`/vasps/${id}`} className="text-body">
                    <span data-testid="name">
                        {row.original.name || 'N/A'}
                    </span>
                    <span className="text-muted font-italic d-block" data-testid="commonName">
                        {row.original.common_name || 'N/A'}
                    </span>
                </Link>
            </p>
        </React.Fragment>
    );
};

export default NameColumn