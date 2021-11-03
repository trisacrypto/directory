import { FetchReviewNotesActionTypes } from "./constants"

const INITIAL_STATE = {
    data: null,
    loading: false
}

const reviewNotesReducers = (state = INITIAL_STATE, action) => {

    switch (action.type) {
        case FetchReviewNotesActionTypes.FETCH_REVIEW_NOTES:
            return {
                ...state,
                loading: true
            }
        case FetchReviewNotesActionTypes.API_RESPONSE_SUCCESS:
            return {
                ...state,
                data: action.payload.data.notes,
                loading: false
            }
        case FetchReviewNotesActionTypes.API_RESPONSE_ERROR:
            return {
                ...state,
                error: action.payload.error
            }
        default:
            return state;
    }
}

export { reviewNotesReducers }