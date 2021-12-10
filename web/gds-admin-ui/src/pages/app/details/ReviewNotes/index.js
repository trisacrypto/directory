import React from 'react';
import { Card } from 'react-bootstrap';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
import { useParams } from 'react-router-dom';
dayjs.extend(relativeTime)

const ReviewNotes = React.lazy(_ => import('./ReviewNotes'))
const ReviewNoteForm = React.lazy(_ => import('./ReviewNoteForm'))


const Comments = () => {
    const params = useParams()

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Reviewer Notes</h4>
                <ReviewNoteForm vaspId={params?.id} />
                <ReviewNotes vaspId={params?.id} />
            </Card.Body>
        </Card>
    );
};

export default Comments;
