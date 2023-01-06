import React from 'react';
import { Button } from 'react-bootstrap';
import { Controller, useForm } from 'react-hook-form';

import TextEditor from '@/components/TextEditor';
import sanitizeMarkdown from '@/utils/sanitize-markdown';
import { useUpdateReviewNote } from '@/features/vasps/services';
import { GET_REVIEW_NOTES } from '@/features/vasps/constants';
import { queryClient } from '@/lib/react-query';
import { toast } from 'react-hot-toast';

function EditReviewNote({ note, handleCancelEditingClick, vaspId, setIsEditable }: any) {
    const { handleSubmit, watch, control, formState } = useForm({
        defaultValues: {
            note: note?.text,
            noteId: note?.id,
        },
    });
    const { mutate: updateReviewNote } = useUpdateReviewNote();
    const watchedNote = watch('note').trim();
    const [isSubmitting, setIsSubmiting] = React.useState(false);

    const onSubmit = (data: any) => {
        const { note, noteId } = data;

        setIsSubmiting(true);
        const sanitizedNote = sanitizeMarkdown(note);

        updateReviewNote(
            {
                vaspId,
                noteId,
                note: sanitizedNote?.trim(),
            },
            {
                onError: (_: any, __: any, context: any) => {
                    if (context?.previousReviewNotes) {
                        queryClient.setQueryData([GET_REVIEW_NOTES], context.previousReviewNotes);
                    }
                    toast.error('Sorry, unable to update review note');
                },
                onSuccess: () => {
                    queryClient.invalidateQueries([GET_REVIEW_NOTES]);
                    toast.success('Review note updated successfully');
                    setIsSubmiting(false);
                    setIsEditable(false);
                },
            }
        );
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <Controller
                name="note"
                control={control}
                render={({ field }) => <TextEditor {...field} className="mb-2" />}
            />
            <div className="d-flex gap-1">
                <Button
                    type="submit"
                    disabled={isSubmitting || !watchedNote || !formState.isDirty}
                    className="btn btn-success btn-sm">
                    Save
                </Button>
                <Button onClick={handleCancelEditingClick} type="button" className="btn btn-sm btn-danger">
                    Cancel
                </Button>
            </div>
        </form>
    );
}

export default EditReviewNote;
