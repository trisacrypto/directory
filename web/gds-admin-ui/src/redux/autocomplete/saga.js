import { call, put, takeEvery } from 'redux-saga/effects';

import getAllAutocomplete from '@/services/autocomplete';

import { FetchAutocompleteActionTypes } from './constants';
import { fetchAutocompleteApiResponseError, fetchAutocompleteApiResponseSuccess } from '.';

function* fetchAutocompletes() {
  try {
    const response = yield call(getAllAutocomplete);
    const {data} = response;
    yield put(fetchAutocompleteApiResponseSuccess(FetchAutocompleteActionTypes.API_RESPONSE_SUCCESS, data));
  } catch (error) {
    console.error(error);
    yield put(fetchAutocompleteApiResponseError(FetchAutocompleteActionTypes.API_RESPONSE_ERROR, error.message));
  }
}

function* autocompletesSaga() {
  yield takeEvery(FetchAutocompleteActionTypes.FETCH_AUTOCOMPLETE, fetchAutocompletes);
}

export default autocompletesSaga;
