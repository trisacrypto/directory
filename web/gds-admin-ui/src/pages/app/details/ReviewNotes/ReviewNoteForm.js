import React from 'react'
import { useForm } from "react-hook-form";
import toast from 'react-hot-toast';
import { useDispatch } from 'react-redux';
import useSafeDispatch from '../../../../hooks/useSafeDispatch';
import { createReviewNoteApiResponseSuccess } from '../../../../redux/review-notes';
import { postReviewNote } from '../../../../services/review-notes';


function ReviewNoteForm({ vaspId }) {
    const { register, handleSubmit, watch, reset } = useForm({
        defaultValues: {
            note: ''
        }
    });
    const [isSubmitting, setIsSubmiting] = React.useState(false)
    const watchedNote = watch('note').trim()
    const dispatch = useDispatch()
    const safeDispatch = useSafeDispatch(dispatch)


    const onSubmit = (data) => {
        const { note } = data
        setIsSubmiting(true)

        postReviewNote(note, vaspId).then(response => {
            const data = response?.data
            console.log('[data]', data)
            if (data) {
                safeDispatch(createReviewNoteApiResponseSuccess(data))
                toast.success('Review note added successfully')
            }
            setIsSubmiting(false)
            reset({ note: '' })
        }).catch(error => {
            console.error('[ReviewNoteForm] onSubmit', error)
            setIsSubmiting(false)
        })
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <textarea
                {...register("note", { required: true })}
                className="form-control form-control-light mb-2"
                placeholder="Write a review note"
                id="review-note"
                rows="3"></textarea>
            <div className="text-end">
                <div className="btn-group mb-2 ms-2">
                    <button type="submit" className="btn btn-primary btn-sm" disabled={!watchedNote || isSubmitting}>
                        Submit
                    </button>
                </div>
            </div>
        </form>
    )
}

export default ReviewNoteForm
