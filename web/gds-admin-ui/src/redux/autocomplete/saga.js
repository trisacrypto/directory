const { call, put, takeEvery } = require("redux-saga/effects")
const { default: getAllAutocomplete } = require("services/autocomplete")
const { fetchAutocompleteApiResponseSuccess, FetchAutocompleteActionTypes, fetchAutocompleteApiResponseError } = require(".")

function* fetchAutocompletes() {
    try {
        const response = yield call(getAllAutocomplete)
        const data = response.data
        yield put(fetchAutocompleteApiResponseSuccess(FetchAutocompleteActionTypes.API_RESPONSE_SUCCESS, data))
    } catch (error) {
        console.error(error)
        yield put(fetchAutocompleteApiResponseError(FetchAutocompleteActionTypes.API_RESPONSE_ERROR, error.message))
    }
}

function* autocompletesSaga() {
    yield takeEvery(FetchAutocompleteActionTypes.FETCH_AUTOCOMPLETE, fetchAutocompletes)
}

export default autocompletesSaga