import { APICore } from '@/helpers/api/apiCore';
import { MutationConfig, queryClient } from '@/lib/react-query';
import { useMutation } from '@tanstack/react-query';
import { GET_REVIEW_NOTES } from '../constants';
import { X_CSRF_TOKEN } from '@/constants';

const apiCore = new APICore();

function updateReviewNote({ note, noteId, vaspId }: any) {
    const data = {
        text: note,
    };
    return apiCore.update(`/vasps/${vaspId}/notes/${noteId}`, data, {
        headers: {
            [X_CSRF_TOKEN]: apiCore.getCsrfToken(),
        },
    });
}

type UseUpdateReviewNotesOptions = {
    config?: MutationConfig<typeof updateReviewNote>;
};

export const useUpdateReviewNote = ({ config }: UseUpdateReviewNotesOptions = {}) => {
    return useMutation({
        ...config,
        mutationFn: updateReviewNote,
        onMutate: async () => {
            await queryClient.cancelQueries([GET_REVIEW_NOTES]);

            const previousReviewNotes = queryClient.getQueryData<any[]>(['get-review-notes']);

            return { previousReviewNotes };
        },
    });
};
