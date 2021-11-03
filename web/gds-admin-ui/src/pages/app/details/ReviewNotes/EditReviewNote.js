import React from 'react'
import { Button } from 'react-bootstrap';
import { useForm } from 'react-hook-form';
import toast from 'react-hot-toast';
import { useDispatch } from 'react-redux';
import useSafeDispatch from '../../../../hooks/useSafeDispatch';
import { updateReviewNoteApiResponseSuccess } from '../../../../redux/review-notes';
import { updateReviewNote } from '../../../../services/review-notes';

function EditReviewNote({ note, handleCancelEditingClick, vaspId, setIsEditable }) {
    const { register, handleSubmit, watch } = useForm({
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

        updateReviewNote(note?.trim(), noteId, vaspId).then(response => {
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
            <textarea {...register("note", { required: true })} name="note" autoFocus className="form-control form-control-light mb-2" rows="3"></textarea>
            <div className='d-flex gap-1'>
                <Button type="submit" disabled={isSubmitting || !watchedNote} className='btn btn-success btn-sm'>Save</Button>
                <Button onClick={handleCancelEditingClick} type="button" className='btn btn-sm btn-danger'>Cancel</Button>
            </div>
        </form>
    )
}

export default EditReviewNote
