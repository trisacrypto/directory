import React from 'react'
import { useForm } from "react-hook-form";
import PropTypes from 'prop-types';


function ReviewNoteForm({ handleReviewNoteSubmit, isSubmitting }) {
    const { register, handleSubmit, watch, reset } = useForm({
        defaultValues: {
            note: ''
        }
    });
    const watchedNote = watch('note').trim()

    const onSubmit = async (data) => {
        const { note } = data
        await handleReviewNoteSubmit(note)
        reset({ note: '' })

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

ReviewNoteForm.propTypes = {
    handleReviewNoteSubmit: PropTypes.func
}

export default ReviewNoteForm
