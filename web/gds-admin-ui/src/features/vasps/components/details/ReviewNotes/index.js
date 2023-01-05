import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import React from 'react';
import { Card } from 'react-bootstrap';
import toast from 'react-hot-toast';
import { useDispatch } from 'react-redux';
import { useParams } from 'react-router-dom';

import useSafeDispatch from '@/hooks/useSafeDispatch';
import { createReviewNoteApiResponseSuccess } from '@/redux/vasp-details';
import { postReviewNote } from '@/services/review-notes';
import sanitizeMarkdown from '@/utils/sanitize-markdown';

dayjs.extend(relativeTime);

const ReviewNotes = React.lazy((_) => import('./ReviewNotes'));
const ReviewNoteForm = React.lazy((_) => import('./ReviewNoteForm'));

const Comments = () => {
  const params = useParams();
  const dispatch = useDispatch();
  const safeDispatch = useSafeDispatch(dispatch);
  const [isSubmitting, setIsSubmiting] = React.useState(false);

  const handleReviewNoteSubmit = (note) => {
    if (note && params && params.id) {
      const sanitizedNote = sanitizeMarkdown(note);
      postReviewNote(sanitizedNote, params.id)
        .then((response) => {
          const data = response?.data;
          safeDispatch(createReviewNoteApiResponseSuccess(data));
          toast.success('Review note added successfully');
        })
        .catch((error) => {
          console.error('[ReviewNoteForm] onSubmit', error);
          setIsSubmiting(false);
        });
    }
  };

  return (
    <Card>
      <Card.Body>
        <h4 className="mt-0 mb-3">Reviewer Notes</h4>
        <ReviewNoteForm
          isSubmitting={isSubmitting}
          handleReviewNoteSubmit={handleReviewNoteSubmit}
        />
        <ReviewNotes vaspId={params?.id} />
      </Card.Body>
    </Card>
  );
};

export default Comments;
