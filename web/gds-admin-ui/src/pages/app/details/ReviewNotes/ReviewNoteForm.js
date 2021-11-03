import React from 'react'

function ReviewNoteForm() {
    return (
        <>
            <textarea
                className="form-control form-control-light mb-2"
                placeholder="Write a review note"
                id="review-note"
                rows="7"></textarea>
            <div className="text-end">
                <div className="btn-group mb-2 ms-2">
                    <button type="button" className="btn btn-primary btn-sm">
                        Submit
                    </button>
                </div>
            </div>
        </>
    )
}

export default ReviewNoteForm
