import React from 'react';
import { Card } from 'react-bootstrap';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
import { useParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { fetchReviewNotesApiResponse } from '../../../../redux/review-notes/actions';
<<<<<<< HEAD
import useSafeDispatch from '../../../../hooks/useSafeDispatch';
=======
>>>>>>> feat: add review notes list
dayjs.extend(relativeTime)

const ReviewNotes = React.lazy(_ => import('./ReviewNotes'))
const ReviewNoteForm = React.lazy(_ => import('./ReviewNoteForm'))


const Comments = () => {
    const params = useParams()
    const dispatch = useDispatch()
<<<<<<< HEAD
    const safeDispatch = useSafeDispatch(dispatch)

    React.useEffect(() => {
        if (params?.id) {
            safeDispatch(fetchReviewNotesApiResponse(params.id))
        }

    }, [params.id, safeDispatch])
=======
    const isMounted = React.useRef(true);


    React.useEffect(() => {
        if (isMounted.current) {
            if (params?.id) {
                dispatch(fetchReviewNotesApiResponse(params.id))
            }
        }

        return () => isMounted.current = false;
    }, [dispatch, params.id])
>>>>>>> feat: add review notes list

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Reviewer Notes</h4>
<<<<<<< HEAD
                <ReviewNoteForm vaspId={params?.id} />
                <ReviewNotes vaspId={params?.id} />
=======
                <ReviewNoteForm />
                <ReviewNotes />
>>>>>>> feat: add review notes list
            </Card.Body>
        </Card>
    );
};

export default Comments;
