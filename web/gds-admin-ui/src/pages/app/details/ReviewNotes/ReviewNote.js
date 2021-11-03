import dayjs from 'dayjs'
import React from 'react'
import Gravatar from 'react-gravatar'
import { generateMd5 } from '../../../../utils'

function ReviewNote({ note, user }) {
    return (
        <div className="d-flex align-items-start mt-2 mb-3">
            <Gravatar default='identicon' md5={generateMd5(user?.email)} protocol="https://" className='me-3 avatar-sm rounded-circle' />
            <div className="w-100 overflow-hidden">
                <div className='d-flex gap-1'>
                    <h5 className="mt-0">{note.author}</h5>
                    <small className='text-muted'>{dayjs(note.created).fromNow()}</small>
                </div>
                {note.text}
            </div>
        </div>
    )
}

export default ReviewNote
