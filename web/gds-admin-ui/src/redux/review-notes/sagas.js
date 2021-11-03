import toast from "react-hot-toast"
import { call, put, takeEvery } from "redux-saga/effects"
import { getReviewNotes } from "../../services/review-notes"
import { fetchReviewNotesApiResponseError, fetchReviewNotesApiResponseSuccess } from "./actions"
import { FetchReviewNotesActionTypes } from "./constants"

function* fetchReviewNotes({ payload: { id } }) {
    try {
        const response = yield call(getReviewNotes, id)
        const data = response?.data
        yield put(fetchReviewNotesApiResponseSuccess(data))
    } catch (error) {
        toast.error(error)
        yield put(fetchReviewNotesApiResponseError(error.message))
    }
}

export function* reviewNotesSaga() {
    yield takeEvery(FetchReviewNotesActionTypes.FETCH_REVIEW_NOTES, fetchReviewNotes)
}