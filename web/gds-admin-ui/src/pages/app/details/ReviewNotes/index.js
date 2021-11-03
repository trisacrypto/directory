import React from 'react';
import { Card } from 'react-bootstrap';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
import { useParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { fetchReviewNotesApiResponse } from '../../../../redux/review-notes/actions';
dayjs.extend(relativeTime)

const ReviewNotes = React.lazy(_ => import('./ReviewNotes'))
const ReviewNoteForm = React.lazy(_ => import('./ReviewNoteForm'))


const Comments = () => {
    const params = useParams()
    const dispatch = useDispatch()
    const isMounted = React.useRef(true);


    React.useEffect(() => {
        if (isMounted.current) {
            if (params?.id) {
                dispatch(fetchReviewNotesApiResponse(params.id))
            }
        }

        return () => isMounted.current = false;
    }, [dispatch, params.id])

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Reviewer Notes</h4>
                <ReviewNoteForm />
                <ReviewNotes />
            </Card.Body>
        </Card>
    );
};

export default Comments;
