import toast from "react-hot-toast"
import { call, put, takeEvery } from "redux-saga/effects"
import { DeleteReviewNotesActionTypes } from "."
import { getReviewNotes, deleteReviewNote as deleteNote } from "../../services/review-notes"
import { fetchReviewNotesApiResponseError, fetchReviewNotesApiResponseSuccess } from "./actions"
import { FetchReviewNotesActionTypes } from "./constants"

function* fetchReviewNotes({ payload: { id } }) {
    try {
        const response = yield call(getReviewNotes, id)
        const { notes } = response?.data
        const sortedData = Array.isArray(notes) ? notes.sort((a, b) => {
            const date1 = a?.modified ? new Date(a?.modified) : new Date(a?.created)
            const date2 = b?.modified ? new Date(b?.modified) : new Date(b?.created)

            return date2 - date1
        }) : []

        yield put(fetchReviewNotesApiResponseSuccess(sortedData))
    } catch (error) {
        toast.error(error)
        yield put(fetchReviewNotesApiResponseError(error.message))
    }
}

function* deleteReviewNote({ payload: { noteId, vaspId } }) {
    try {
        const response = yield call(deleteNote, noteId, vaspId)
        if (response) {
            toast.success('The note has been deleted successfully')
        }
    } catch (error) {
        toast.error(error.message)
    }
}

export function* reviewNotesSaga() {
    yield takeEvery(FetchReviewNotesActionTypes.FETCH_REVIEW_NOTES, fetchReviewNotes)
}

export function* deleteReviewNoteSaga() {
    yield takeEvery(DeleteReviewNotesActionTypes.DELETE_REVIEW_NOTES, deleteReviewNote)
}
