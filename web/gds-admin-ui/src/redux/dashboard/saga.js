
import axios from "axios"
import { call, put, takeEvery, fork, all } from "redux-saga/effects"
import { getSummary, getVasps } from "../../services/dashboard"
import { fetchVaspsApiResponseSuccess, fetchVaspsApiResponseError, fetchSummaryApiResponseSuccess, fetchSummaryApiResponseError, fetchRegistrationsReviewsSuccess, fetchRegistrationsReviewsError } from "./actions"
import { FetchSummaryActionTypes, FetchVaspsActionTypes } from "./constants"



function* fetchSummary() {
    try {
        const response = yield call(getSummary)
        const data = response.data
        yield put(fetchSummaryApiResponseSuccess(FetchVaspsActionTypes.API_RESPONSE_SUCCESS, data))
    } catch (error) {
        yield put(fetchSummaryApiResponseError(FetchVaspsActionTypes.API_RESPONSE_ERROR, error.message))
    }
}


function* fetchVasps() {
    try {
        const response = yield call(getVasps, { status: 'pending+review' })
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


function* dashboardSaga() {
    yield all([
        fork(summarySaga),
        fork(vaspsSaga),
    ])
}

export default dashboardSaga;