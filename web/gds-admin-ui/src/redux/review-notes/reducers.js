import { DeleteReviewNotesActionTypes } from "."
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
                data: action.payload.data,
                loading: false
            }
        case FetchReviewNotesActionTypes.API_RESPONSE_ERROR:
            return {
                ...state,
                error: action.payload.error
            }
        case DeleteReviewNotesActionTypes.DELETE_REVIEW_NOTES:
            const filteredData = state.data.filter(note => note && note.id !== action.payload.noteId)
            return {
                ...state,
                data: filteredData
            }
        default:
            return state;
    }
}

export { reviewNotesReducers }