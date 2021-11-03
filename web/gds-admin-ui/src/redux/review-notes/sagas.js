import toast from "react-hot-toast"
import { call, put, takeEvery } from "redux-saga/effects"
import { getReviewNotes } from "../../services/review-notes"
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

export function* reviewNotesSaga() {
    yield takeEvery(FetchReviewNotesActionTypes.FETCH_REVIEW_NOTES, fetchReviewNotes)
}