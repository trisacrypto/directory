import { APICore } from '@/helpers/api/apiCore';
import { MutationConfig, queryClient } from '@/lib/react-query';
import { useMutation } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import { GET_REVIEW_NOTES } from '../constants';
import { X_CSRF_TOKEN } from '@/constants';

const apiCore = new APICore();

function removeReviewNote({ noteId, vaspId }: { noteId: string; vaspId: string }) {
    return apiCore.delete(`/vasps/${vaspId}/notes/${noteId}`, {
        headers: {
            [X_CSRF_TOKEN]: apiCore.getCsrfToken(),
        },
    });
}

type UseGetReviewNotesOptions = {
    config?: MutationConfig<typeof removeReviewNote>;
    noteId: string;
};

export const useDeleteReviewNote = ({ config, noteId }: UseGetReviewNotesOptions) => {
    return useMutation({
        onMutate: async (deletedNote: any) => {
            await queryClient.cancelQueries([GET_REVIEW_NOTES]);

            const previousReviewNotes = queryClient.getQueryData<any[]>(['get-review-notes']);

            queryClient.setQueryData(
                [GET_REVIEW_NOTES],
                previousReviewNotes?.filter((note: any) => note.id !== deletedNote.noteId)
            );

            return { previousReviewNotes };
        },
        onError: (_: any, __: any, context: any) => {
            if (context?.previousReviewNotes) {
                queryClient.setQueryData([GET_REVIEW_NOTES], context.previousReviewNotes);
            }
            toast.error('Sorry, unable to delete review note. Try later');
        },
        onSuccess: () => {
            queryClient.invalidateQueries([GET_REVIEW_NOTES]);
            toast.success('Review note deleted successfully');
        },
        ...config,
        mutationFn: removeReviewNote,
    });
};
