import React from "react"
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(relativeTime)

const LastUpdatedColumn = ({ row }) => {

    return <React.Fragment>
        <p className="m-0 d-inline-block align-middle font-16">
            <span data-testid="last_updated">
                {row?.original?.last_updated ? dayjs(row?.original?.last_updated).fromNow() : 'N/A'}
            </span>
        </p>
    </React.Fragment>
}

export default LastUpdatedColumn