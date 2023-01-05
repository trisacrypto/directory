import classNames from 'classnames';
import React from 'react';

import { StatusLabel } from '@/constants';
import { getStatusClassName } from '@/utils';

const StatusColumn = ({ row }) => (
  <span
    data-testid="verification_status"
    className={classNames('badge', getStatusClassName(row.original.verification_status))}
  >
    {StatusLabel[row.original.verification_status] || 'N/A'}
  </span>
);

export default StatusColumn;
