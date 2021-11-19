import { CreateReviewNoteActionTypes, DeleteReviewNotesActionTypes, UpdateReviewNotesActionTypes } from "."
import { FetchReviewNotesActionTypes } from "./constants"
import { produce } from 'immer'

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
        case UpdateReviewNotesActionTypes.API_RESPONSE_SUCCESS:
            return produce(state, draft => {
                const idx = draft.data.findIndex(note => note?.id === action.payload.note?.id)

                if (idx !== -1) {
                    draft.data[idx] = action.payload.note
                }
            })
        case CreateReviewNoteActionTypes.API_RESPONSE_SUCCESS:
            return produce(state, draft => {
                draft.data.unshift(action.payload.note)
            })
        default:
            return state;
    }
}

export { reviewNotesReducers }