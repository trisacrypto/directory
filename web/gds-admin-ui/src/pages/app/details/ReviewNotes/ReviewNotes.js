import React from 'react'
import dayjs from 'dayjs';
import ReviewNote from './ReviewNote';
import useUser from '../../../../hooks/useUser';
import relativeTime from 'dayjs/plugin/relativeTime'
import { useSelector } from 'react-redux';
import { getAllReviewNotes, getReviewNoteLoadingState } from '../../../../redux/selectors';
import Loader from '../../../../components/Loader';
dayjs.extend(relativeTime)


function ReviewNotes({ data }) {
    const { user } = useUser()
    const reviewNotes = useSelector(getAllReviewNotes)
    const isLoading = useSelector(getReviewNoteLoadingState)

    if (isLoading) return <Loader width={50} />
    if (reviewNotes && reviewNotes.length) {
        return data && data.map((note) => (
            <ReviewNote user={user?.email} note={note} key={note.id} />
        ))
    }

    return <div className='text-center fst-italic text-muted'>No reviewer notes</div>
}

export default ReviewNotes
