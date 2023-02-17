import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

import { useGetReviewNotes } from '@/features/vasps/services';
import { useParams } from 'react-router-dom';
import ReviewNote from './ReviewNote';

dayjs.extend(relativeTime);

function ReviewNotes() {
    const params = useParams<{ id: string }>();
    const { data: reviewNotes } = useGetReviewNotes({ vaspId: params?.id });

    if (reviewNotes && reviewNotes.length < 1) {
        return <div className="text-center fst-italic text-muted">No reviewer notes</div>;
    }

    return (reviewNotes || []).map((note: any) => note.id && <ReviewNote note={note} key={note.id} />);
}

export default ReviewNotes;
