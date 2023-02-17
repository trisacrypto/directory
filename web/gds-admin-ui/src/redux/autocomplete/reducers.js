const { FetchAutocompleteActionTypes } = require('./constants');

const INITIAL_STATE = {
  data: [],
  loading: false,
};

const autocompleteReducers = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case FetchAutocompleteActionTypes.FETCH_AUTOCOMPLETE:
      return {
        ...state,
        loading: true,
      };
    case FetchAutocompleteActionTypes.API_RESPONSE_SUCCESS:
      return {
        ...state,
        data: action.payload.data,
        loading: false,
      };
    case FetchAutocompleteActionTypes.API_RESPONSE_ERROR:
      return {
        ...state,
        error: action.payload.error,
      };
    default:
      return state;
  }
};

export { autocompleteReducers };
