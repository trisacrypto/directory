import { APICore } from '@/helpers/api/apiCore';
import { MutationConfig, queryClient } from '@/lib/react-query';
import { useMutation } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import { GET_REVIEW_NOTES } from '../constants';
import { X_CSRF_TOKEN } from '@/constants';

const apiCore = new APICore();

function postReviewNote({ note, vaspId }: { note: any; vaspId: string }) {
    const payload = { text: note, note_id: '' };
    return apiCore.create(`/vasps/${vaspId}/notes`, payload, {
        headers: {
            [X_CSRF_TOKEN]: apiCore.getCsrfToken(),
        },
    });
}
type UseUpdateReviewNotesOptions = {
    config?: MutationConfig<typeof postReviewNote>;
};

export const useCreateReviewNote = ({ config }: UseUpdateReviewNotesOptions = {}) => {
    return useMutation({
        onError: (_: any, __: any, context: any) => {
            toast.error('Sorry, unable to create a review note. Try later');
        },
        onSuccess: (data) => {
            queryClient.invalidateQueries([GET_REVIEW_NOTES]);
            toast.success('Review note created successfully');
        },
        ...config,
        mutationFn: postReviewNote,
    });
};
