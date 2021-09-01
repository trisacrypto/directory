
import axios from "axios"
import { call, put, takeEvery, fork, all } from "redux-saga/effects"
import { APICore } from "../../helpers/api/apiCore"
import { getSummary } from "../../services/dashboard"
import { fetchCertificateApiResponseError, fetchCertificateApiResponseSuccess, fetchVaspsApiResponseSuccess, fetchVaspsApiResponseError, fetchSummaryApiResponseSuccess, fetchSummaryApiResponseError } from "./actions"
import { FetchCertificatesActionTypes, FetchSummaryActionTypes, FetchVaspsActionTypes } from "./constants"

const api = new APICore()


function* fetchSummary() {
    try {
        const response = yield call(getSummary)
        const data = response.data
        yield put(fetchSummaryApiResponseSuccess(FetchVaspsActionTypes.API_RESPONSE_SUCCESS, data))
    } catch (error) {
        yield put(fetchSummaryApiResponseError(FetchVaspsActionTypes.API_RESPONSE_ERROR, error.message))
    }
}

function* fetchCertificates() {
    try {
        const response = yield call(api.get, "/certificates")
        const data = response && response.data
        yield put(fetchCertificateApiResponseSuccess(FetchCertificatesActionTypes.API_RESPONSE_SUCCESS, data))
    } catch (error) {
        yield put(fetchCertificateApiResponseError(FetchCertificatesActionTypes.API_RESPONSE_ERROR, error.message))
    }
}

function* fetchVasps() {
    try {
        const response = yield call(axios.get, "/vasps")
        const data = response.data
        yield put(fetchVaspsApiResponseSuccess(FetchVaspsActionTypes.API_RESPONSE_SUCCESS, data))
    } catch (error) {
        yield put(fetchVaspsApiResponseError(FetchVaspsActionTypes.API_RESPONSE_ERROR, error.message))
    }
}

export function* summarySaga() {
    yield takeEvery(FetchSummaryActionTypes.FETCH_SUMMARY, fetchSummary)
}

export function* vaspsSaga() {
    yield takeEvery([FetchVaspsActionTypes.FETCH_VASPS], fetchVasps);
}

export function* certificatesSaga() {
    yield takeEvery([FetchCertificatesActionTypes.FETCH_CERTIFICATES], fetchCertificates);
}

function* dashboardSaga() {
    yield all([
        fork(summarySaga),
        fork(vaspsSaga),
        fork(certificatesSaga)
    ])
}

export default dashboardSaga;