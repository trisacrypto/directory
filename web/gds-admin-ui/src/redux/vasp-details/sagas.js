import NProgress from 'nprogress'
import { call, put, takeEvery } from 'redux-saga/effects'
import { getReviewNotes, deleteReviewNote as deleteNote } from 'services/review-notes'
import updateTrixoForm from 'services/trixo'
import { getVasp } from 'services/vasp'
import { fetchVaspDetailsApiResponseSuccess, DeleteReviewNotesActionTypes, fetchReviewNotesApiResponseError, fetchReviewNotesApiResponseSuccess, FetchVaspDetailsActionTypes, fetchVaspDetailsApiResponseError, UpdateTrixoActionTypes, updateTrixoResponseSuccess, updateTrixoResponseError } from '.'


function* fetchVaspDetails({ payload: { id, history } }) {
    NProgress.start()
    try {
        const response = yield call(getVasp, id)
        yield put(fetchVaspDetailsApiResponseSuccess(response.data))
        NProgress.done()
    } catch (error) {
        yield put(fetchVaspDetailsApiResponseError(error.message))
        history.push('/not-found', { error: "Could not retrieve VASP record by ID" })
        NProgress.done()
    }
}

function* updateTrixo({ payload: { trixo, id, setIsOpen } }) {
    try {
        const response = yield call(updateTrixoForm, id, { trixo })
        yield put(updateTrixoResponseSuccess(response.data))

        setIsOpen(false)
    } catch (error) {
        yield put(updateTrixoResponseError(error.message))
        console.error('[updateVaspDetails] error', error.message)
    }
}

function* fetchReviewNotes({ payload: { id } }) {
    NProgress.start()
    try {
        const response = yield call(getReviewNotes, id)
        const { notes } = response.data
        const sortedData = Array.isArray(notes) ? notes.sort((a, b) => {
            const date1 = a?.modified ? new Date(a?.modified) : new Date(a?.created)
            const date2 = b?.modified ? new Date(b?.modified) : new Date(b?.created)

            return date2 - date1
        }) : []

        yield put(fetchReviewNotesApiResponseSuccess(sortedData))
        NProgress.done()
    } catch (error) {
        yield put(fetchReviewNotesApiResponseError(error.message))
        NProgress.done()
    }
}

function* deleteReviewNote({ payload: { noteId, vaspId } }) {
    try {
        const response = yield call(deleteNote, noteId, vaspId)
        if (response) {
            NProgress.done()
        }
    } catch (error) {
        console.error(error)
    }
}

function* vaspDetailsSaga() {
    yield takeEvery(FetchVaspDetailsActionTypes.FETCH_VASP_DETAILS, fetchVaspDetails)
}

export function* reviewNotesSaga() {
    yield takeEvery(FetchVaspDetailsActionTypes.FETCH_VASP_DETAILS, fetchReviewNotes)
}

export function* deleteReviewNoteSaga() {
    yield takeEvery(DeleteReviewNotesActionTypes.DELETE_REVIEW_NOTES, deleteReviewNote)
}

export function* updateTrixoSaga() {
    yield takeEvery(UpdateTrixoActionTypes.UPDATE_TRIXO, updateTrixo)
}

export { vaspDetailsSaga }