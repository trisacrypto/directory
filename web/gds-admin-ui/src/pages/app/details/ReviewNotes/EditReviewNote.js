import React from 'react'
import { Button } from 'react-bootstrap';
import { useForm, Controller } from 'react-hook-form';
import toast from 'react-hot-toast';
import { useDispatch } from 'react-redux';
import useSafeDispatch from 'hooks/useSafeDispatch';
import { updateReviewNote } from 'services/review-notes';
import PropTypes from 'prop-types';
import { updateReviewNoteApiResponseSuccess } from 'redux/vasp-details';
import TextEditor from 'components/TextEditor';
import sanitizeMarkdown from 'utils/sanitize-markdown';

function EditReviewNote({ note, handleCancelEditingClick, vaspId, setIsEditable }) {
    const { register, handleSubmit, watch, control } = useForm({
        defaultValues: {
            note: note?.text
        }
    });
    const watchedNote = watch('note').trim()
    const [isSubmitting, setIsSubmiting] = React.useState(false)
    const dispatch = useDispatch()
    const safeDispatch = useSafeDispatch(dispatch)


    const onSubmit = (data) => {
        const { note, noteId } = data
        setIsSubmiting(true)
        const sanitizedNote = sanitizeMarkdown(note)

        updateReviewNote(sanitizedNote?.trim(), noteId, vaspId).then(response => {
            const note = response?.data
            if (note) {
                safeDispatch(updateReviewNoteApiResponseSuccess(note))

                setIsEditable(false)
                toast.success('The note has been updated successfully')
            } else {
                toast.error('Sorry, you are not able to update this note')
            }
            setIsSubmiting(false)
        }).catch(err => {
            console.error('[EditReviewNote]', err)
            toast.error('Unable to delete this note')
            setIsSubmiting(false)
        })
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <input {...register("noteId")} type="text" defaultValue={note?.id} hidden />
            <Controller name='note' control={control} render={({ field }) => (
                <TextEditor {...field} className="mb-2" />
            )} />
            <div className='d-flex gap-1'>
                <Button type="submit" disabled={isSubmitting || !watchedNote} className='btn btn-success btn-sm'>Save</Button>
                <Button onClick={handleCancelEditingClick} type="button" className='btn btn-sm btn-danger'>Cancel</Button>
            </div>
        </form>
    )
}

EditReviewNote.propTypes = {
    note: PropTypes.object.isRequired,
    handleCancelEditingClick: PropTypes.func.isRequired,
    vaspId: PropTypes.string.isRequired,
    setIsEditable: PropTypes.func.isRequired
}

export default EditReviewNote
