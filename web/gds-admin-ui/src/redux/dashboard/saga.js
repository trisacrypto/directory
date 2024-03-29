import NProgress from 'nprogress';
import toast from 'react-hot-toast';
import { all, call, fork, put, takeEvery } from 'redux-saga/effects';

import { getRegistrationReviews, getSummary, getVasps } from '@/services/dashboard';

import {
  fetchRegistrationsReviewsError,
  fetchRegistrationsReviewsSuccess,
  fetchSummaryApiResponseError,
  fetchSummaryApiResponseSuccess,
  fetchVaspsApiResponseError,
  fetchVaspsApiResponseSuccess,
} from './actions';
import {
  FetchPendingVaspsActionTypes,
  FetchRegistrationsReviewsActionTypes,
  FetchSummaryActionTypes,
  FetchVaspsActionTypes,
} from './constants';

function* fetchSummary() {
  NProgress.start();
  try {
    const response = yield call(getSummary);
    const data = response?.data;
    yield put(fetchSummaryApiResponseSuccess(FetchVaspsActionTypes.API_RESPONSE_SUCCESS, data));
    NProgress.done();
  } catch (error) {
    toast.error(error);
    yield put(
      fetchSummaryApiResponseError(FetchVaspsActionTypes.API_RESPONSE_ERROR, error.message)
    );
    NProgress.done();
  }
}

function* fetchPendingVasps() {
  NProgress.start();
  try {
    const response = yield call(getVasps, 'status=pending+review');
    const { data } = response;
    yield put(fetchVaspsApiResponseSuccess(FetchVaspsActionTypes.API_RESPONSE_SUCCESS, data));
    NProgress.done();
  } catch (error) {
    yield put(fetchVaspsApiResponseError(error.message));
    NProgress.done();
  }
}

function* fetchVasps({ payload }) {
  NProgress.start();
  try {
    const response = yield call(getVasps, payload?.queryParams);
    const data = response?.data;

    yield put(fetchVaspsApiResponseSuccess(FetchVaspsActionTypes.API_RESPONSE_SUCCESS, data));
    NProgress.done();
  } catch (error) {
    toast.error(error);
    yield put(fetchVaspsApiResponseError(FetchVaspsActionTypes.API_RESPONSE_ERROR, error.message));
    NProgress.done();
  }
}

function* fecthRegistrationsReviews() {
  NProgress.start();
  try {
    const response = yield call(getRegistrationReviews);
    const data = response?.data;

    yield put(
      fetchRegistrationsReviewsSuccess(
        FetchRegistrationsReviewsActionTypes.API_RESPONSE_SUCCESS,
        data
      )
    );
    NProgress.done();
  } catch (error) {
    toast.error(error);
    yield put(
      fetchRegistrationsReviewsError(
        FetchRegistrationsReviewsActionTypes.API_RESPONSE_ERROR,
        error.message
      )
    );
    NProgress.done();
  }
}

export function* summarySaga() {
  yield takeEvery(FetchSummaryActionTypes.FETCH_SUMMARY, fetchSummary);
}

export function* vaspsSaga() {
  yield takeEvery([FetchVaspsActionTypes.FETCH_VASPS], fetchVasps);
}

export function* pendingVaspsSaga() {
  yield takeEvery(FetchPendingVaspsActionTypes.FETCH_PENDING_VASPS, fetchPendingVasps);
}

export function* registrationReviews() {
  yield takeEvery(
    FetchRegistrationsReviewsActionTypes.FETCH_REGISTRATIONS_REVIEWS,
    fecthRegistrationsReviews
  );
}

function* dashboardSaga() {
  yield all([fork(summarySaga), fork(pendingVaspsSaga), fork(registrationReviews)]);
}

export default dashboardSaga;
