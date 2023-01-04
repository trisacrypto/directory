import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import React from 'react';

dayjs.extend(relativeTime);

const LastUpdatedColumn = ({ row }) => (
  <p className="m-0 d-inline-block align-middle font-16">
    <span data-testid="last_updated">
      {row?.original?.last_updated ? dayjs(row?.original?.last_updated).fromNow() : 'N/A'}
    </span>
  </p>
);

export default LastUpdatedColumn;
