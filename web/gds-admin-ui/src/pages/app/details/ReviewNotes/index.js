import React from 'react';
import { Card } from 'react-bootstrap';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
import { useParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import useSafeDispatch from 'hooks/useSafeDispatch';
import toast from 'react-hot-toast';
import { postReviewNote } from 'services/review-notes';
import { createReviewNoteApiResponseSuccess } from 'redux/vasp-details';
dayjs.extend(relativeTime)

const ReviewNotes = React.lazy(_ => import('./ReviewNotes'))
const ReviewNoteForm = React.lazy(_ => import('./ReviewNoteForm'))


const Comments = () => {
    const params = useParams()
    const dispatch = useDispatch()
    const safeDispatch = useSafeDispatch(dispatch)
    const [isSubmitting, setIsSubmiting] = React.useState(false)

    const handleReviewNoteSubmit = (note) => {
        if (note && (params && params.id)) {
            postReviewNote(note, params.id).then(response => {
                const data = response?.data
                safeDispatch(createReviewNoteApiResponseSuccess(data))
                toast.success('Review note added successfully')
            }).catch(error => {
                console.error('[ReviewNoteForm] onSubmit', error)
                setIsSubmiting(false)
            })
        }
    }

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Reviewer Notes</h4>
                <ReviewNoteForm isSubmitting={isSubmitting} handleReviewNoteSubmit={handleReviewNoteSubmit} />
                <ReviewNotes vaspId={params?.id} />
            </Card.Body>
        </Card>
    );
};

export default Comments;
