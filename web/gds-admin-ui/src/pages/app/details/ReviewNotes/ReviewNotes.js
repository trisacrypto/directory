import React from 'react'
import dayjs from 'dayjs';
import ReviewNote from './ReviewNote';
import useUser from 'hooks/useUser';
import relativeTime from 'dayjs/plugin/relativeTime'
import { useSelector } from 'react-redux';
import { getAllReviewNotes, getReviewNoteLoadingState } from 'redux/selectors';
import Loader from 'components/Loader';
import PropTypes from 'prop-types';
dayjs.extend(relativeTime)


function ReviewNotes({ vaspId }) {
    const { user } = useUser()
    const reviewNotes = useSelector(getAllReviewNotes)
    const isLoading = useSelector(getReviewNoteLoadingState)

    if (isLoading) {
        return <Loader width={50} />
    }
    if (!isLoading && reviewNotes?.length) {
        return reviewNotes.map((note) => note.id && (
            <ReviewNote vaspId={vaspId} user={user?.email} note={note} key={note.id} />
        ))
    }

    return !reviewNotes?.length && <div className='text-center fst-italic text-muted'>No reviewer notes</div>
}

ReviewNotes.propTypes = {
    vaspId: PropTypes.string.isRequired
}

export default ReviewNotes
