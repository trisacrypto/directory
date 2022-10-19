import React from 'react'
import { useForm } from "react-hook-form";
import PropTypes from 'prop-types';
import TextEditor from 'components/TextEditor';
import localforage from 'localforage'
import { useParams } from 'react-router-dom'
import { isEditorContentEmpty } from 'utils'


function ReviewNoteForm({ handleReviewNoteSubmit, isSubmitting }) {
    const [note, setNote] = React.useState({})
    const params = useParams()
    const vaspId = params?.id || ''
    const storageKey = `review_note_draft/${vaspId}`
    const { handleSubmit, reset } = useForm({
        defaultValues: {
            note,
        }
    });

    React.useEffect(() => {
        (async () => {
            const reviewNoteDraft = await localforage.getItem(storageKey)
            setNote(reviewNoteDraft)
        })()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [vaspId])

    React.useEffect(() => {
        localforage.setItem(storageKey, { ...note, vaspId })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [note, vaspId])


    const onSubmit = async (note) => {
        await handleReviewNoteSubmit(note)
        reset()
        localforage.removeItem(storageKey)
        setNote('')
    }


    const handleTextChange = (html) => setNote({ ...note, reviewNote: html })

    return (
        <form onSubmit={handleSubmit(() => onSubmit(note.reviewNote))}>
            <TextEditor value={note?.reviewNote} className="mb-2" onChange={handleTextChange} data-testid="text-editor" />
            <div className="text-end">
                <div className="btn-group mb-2 ms-2">
                    <button type="submit" className="btn btn-primary btn-sm" disabled={isEditorContentEmpty(note?.reviewNote) || isSubmitting}>
                        Submit
                    </button>
                </div>
            </div>
        </form >
    )
}

ReviewNoteForm.propTypes = {
    handleReviewNoteSubmit: PropTypes.func
}

export default ReviewNoteForm
