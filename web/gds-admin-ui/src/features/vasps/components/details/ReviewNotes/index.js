import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import React from 'react';
import { Card } from 'react-bootstrap';
import { useParams } from 'react-router-dom';

import sanitizeMarkdown from '@/utils/sanitize-markdown';
import { useCreateReviewNote } from '@/features/vasps/services';

dayjs.extend(relativeTime);

const ReviewNotes = React.lazy((_) => import('./ReviewNotes'));
const ReviewNoteForm = React.lazy((_) => import('./ReviewNoteForm'));

const Comments = () => {
    const params = useParams();
    const [isSubmitting, setIsSubmiting] = React.useState(false);
    const { mutate: createReviewNote, isSuccess } = useCreateReviewNote();

    const handleReviewNoteSubmit = (note) => {
        if (note && params && params.id) {
            const sanitizedNote = sanitizeMarkdown(note);

            createReviewNote({
                note: sanitizedNote,
                vaspId: params?.id,
            });

            if (isSuccess) {
                setIsSubmiting(false);
            }
        }
    };

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Reviewer Notes</h4>
                <ReviewNoteForm isSubmitting={isSubmitting} handleReviewNoteSubmit={handleReviewNoteSubmit} />
                <ReviewNotes />
            </Card.Body>
        </Card>
    );
};

export default Comments;
