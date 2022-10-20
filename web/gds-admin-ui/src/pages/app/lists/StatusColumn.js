import React from "react";
import classNames from 'classnames';
import { getStatusClassName } from 'utils';
import { StatusLabel } from '../../../constants';


const StatusColumn = ({ row }) => {
    return (
        <React.Fragment>
            <span
                data-testid="verification_status"
                className={classNames('badge', getStatusClassName(row.original.verification_status))}>
                {StatusLabel[row.original.verification_status] || 'N/A'}
            </span>
        </React.Fragment>
    );
};

export default StatusColumn