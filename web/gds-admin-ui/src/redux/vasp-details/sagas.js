import NProgress from 'nprogress'
import { call, put, takeEvery } from 'redux-saga/effects'
import { getReviewNotes, deleteReviewNote as deleteNote } from 'services/review-notes'
import updateTrixoForm from 'services/trixo'
import { getVasp, updateVasp } from 'services/vasp'
import { fetchVaspDetailsApiResponseSuccess, DeleteReviewNotesActionTypes, fetchReviewNotesApiResponseError, fetchReviewNotesApiResponseSuccess, FetchVaspDetailsActionTypes, fetchVaspDetailsApiResponseError, UpdateTrixoActionTypes, updateTrixoResponseSuccess, updateTrixoResponseError, updateBusinessInfosResponseSuccess, updateBusinessInfosResponseError, UpdateBusinessInfosActionTypes, updateTrisaImplementationDetailsResponseSuccess, updateTrisaImplementationDetailsResponseError, UpdateTrisaImplementationDetailsActionTypes, updateIvms101ResponseError, UpdateIvms101ActionTypes, updateIvms101ResponseSuccess } from '.'


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

function* updateBusinessInfos({ payload: { businessInfos, id, setIsOpen } }) {
    try {
        const response = yield call(updateVasp, id, businessInfos)
        yield put(updateBusinessInfosResponseSuccess(response.data))

        setIsOpen(false)
    } catch (error) {
        yield put(updateBusinessInfosResponseError(error.message))
        console.error('[updateBusinessInfos] error', error.message)
    }
}

function* updateTrisaDetails({ payload: { trisa, id, setIsOpen } }) {
    try {
        const response = yield call(updateVasp, id, trisa)
        yield put(updateTrisaImplementationDetailsResponseSuccess(response.data))

        setIsOpen(false)
    } catch (error) {
        const errorMessage = error.response && error.response.data ? { message: error.response.data['error'], errorStatus: error.response['status'], statusText: error.response['statusText'] } : 'Something went wrong'
        yield put(updateTrisaImplementationDetailsResponseError(errorMessage))
    }
}

function* updateIvms({ payload: { ivms, id, setIsOpen } }) {
    try {
        const response = yield call(updateVasp, id, ivms)
        yield put(updateIvms101ResponseSuccess(response.data))

        setIsOpen(false)
    } catch (error) {
        const errorMessage = error.response && error.response.data ? { message: error.response.data['error'], errorStatus: error.response['status'], statusText: error.response['statusText'] } : 'Something went wrong'
        yield put(updateIvms101ResponseError(errorMessage))
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

export function* updateBusinessInfosSaga() {
    yield takeEvery(UpdateBusinessInfosActionTypes.UPDATE_BUSINESS_INFOS, updateBusinessInfos)
}

export function* updateTrisaDetailsSaga() {
    yield takeEvery(UpdateTrisaImplementationDetailsActionTypes.UPDATE_TRISA_DETAILS, updateTrisaDetails)
}

export function* updateIvmsSaga() {
    yield takeEvery(UpdateIvms101ActionTypes.UPDATE_IVMS_101, updateIvms)
}

export { vaspDetailsSaga }