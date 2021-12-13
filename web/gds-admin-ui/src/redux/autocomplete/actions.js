const { FetchAutocompleteActionTypes } = require("./constants");

const fetchAutocompleteApiResponseSuccess = (actionType, data) => ({
    type: FetchAutocompleteActionTypes.API_RESPONSE_SUCCESS,
    payload: { actionType, data },
});

const fetchAutocompleteApiResponseError = (actionType, data) => ({
    type: FetchAutocompleteActionTypes.API_RESPONSE_ERROR,
    payload: { actionType, data },
});

const fetchAutocomplete = () => ({
    type: FetchAutocompleteActionTypes.FETCH_AUTOCOMPLETE,
    payload: {},
});

export { fetchAutocomplete, fetchAutocompleteApiResponseError, fetchAutocompleteApiResponseSuccess }