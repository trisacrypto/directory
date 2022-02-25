import { call, put, takeEvery } from "redux-saga/effects"
import getAllAutocomplete from "services/autocomplete"
import { fetchAutocompleteApiResponseError, fetchAutocompleteApiResponseSuccess } from "."
import { FetchAutocompleteActionTypes } from "./constants"

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