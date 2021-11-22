import dayjs from 'dayjs'
import React from 'react'
import Gravatar from 'react-gravatar'
import { useDispatch } from 'react-redux'
import useSafeDispatch from 'hooks/useSafeDispatch'
import { deleteReviewNoteApiResponse } from 'redux/review-notes'
import { generateMd5 } from 'utils'
import EditReviewNote from './EditReviewNote'

function ReviewNote({ note, vaspId }) {
    const [isEditable, setIsEditable] = React.useState(false);
    const dispatch = useDispatch()
    const safeDispatch = useSafeDispatch(dispatch)

    const handleDeleteClick = () => {
        if (note && vaspId) {
            if (window.confirm("Are you sure you want to delete this reviewer's note?")) {
                safeDispatch(deleteReviewNoteApiResponse(note?.id, vaspId))
            }
        }
    }

    const handleEditClick = () => {
        setIsEditable(true)
    }

    const handleCancelEditingClick = () => setIsEditable(!isEditable)

    return (
        <div className='d-flex align-items-start mt-2 mb-3'>
            <Gravatar default='identicon' md5={generateMd5(note?.author)} protocol="https://" className='me-3 avatar-sm rounded-circle' />
            <div className="w-100 overflow-hidden">
                <div className='d-flex justify-content-between pt-1'>
                    <div className='d-flex flex-column'>
                        <h5 className="m-0">{note.author}</h5>
                        {
                            note?.modified ?
                                <small className='text-muted d-block fst-italic mb-1'>edited {dayjs(note.updated).fromNow()}</small> :
                                <small className='text-muted d-block fst-italic mb-1'>created {dayjs(note.created).fromNow()}</small>
                        }
                    </div>
                    <div hidden={isEditable}>
                        <button onClick={handleEditClick} className='py-0 px-1 btn btn-success me-sm-1 me-xl-1'><i className='mdi mdi-square-edit-outline'></i> <small className='d-xs-none d-sm-none d-md-none d-lg-none d-xl-inline'>Edit</small></button>
                        <button onClick={handleDeleteClick} className='py-0 px-1 btn btn-danger'><i className='mdi mdi-trash-can-outline'></i> <small className='d-xs-none d-sm-none d-md-none d-lg-none d-xl-inline'>Delete</small></button>
                    </div>
                </div>
                <div>
                    {isEditable ? <EditReviewNote note={note} vaspId={vaspId} setIsEditable={setIsEditable} handleCancelEditingClick={handleCancelEditingClick} /> : <>{note.text}</>}
                </div>
            </div>
        </div>
    )
}

export default ReviewNote
