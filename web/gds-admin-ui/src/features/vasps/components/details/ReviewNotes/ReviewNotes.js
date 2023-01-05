import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import PropTypes from 'prop-types';
import { useSelector } from 'react-redux';

import Loader from '@/components/Loader';
import useUser from '@/hooks/useUser';
import { getAllReviewNotes, getReviewNoteLoadingState } from '@/redux/selectors';

import ReviewNote from './ReviewNote';

dayjs.extend(relativeTime);

function ReviewNotes({ vaspId }) {
  const { user } = useUser();
  const reviewNotes = useSelector(getAllReviewNotes);
  const isLoading = useSelector(getReviewNoteLoadingState);

  if (isLoading) {
    return <Loader width={50} />;
  }
  if (!isLoading && reviewNotes?.length) {
    return reviewNotes.map(
      (note) =>
        note.id && <ReviewNote vaspId={vaspId} user={user?.email} note={note} key={note.id} />
    );
  }

  return (
    !reviewNotes?.length && (
      <div className="text-center fst-italic text-muted">No reviewer notes</div>
    )
  );
}

ReviewNotes.propTypes = {
  vaspId: PropTypes.string.isRequired,
};

export default ReviewNotes;
