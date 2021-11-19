import { CreateReviewNoteActionTypes, DeleteReviewNotesActionTypes, UpdateReviewNotesActionTypes } from ".";
import { FetchReviewNotesActionTypes } from "./constants";

const fetchReviewNotesApiResponse = (vaspId) => ({
    type: FetchReviewNotesActionTypes.FETCH_REVIEW_NOTES,
    payload: { id: vaspId },
});

const fetchReviewNotesApiResponseSuccess = (data) => ({
    type: FetchReviewNotesActionTypes.API_RESPONSE_SUCCESS,
    payload: { data },
});

const fetchReviewNotesApiResponseError = (error) => ({
    type: FetchReviewNotesActionTypes.API_RESPONSE_ERROR,
    payload: { error },
});

const deleteReviewNoteApiResponse = (noteId, vaspId) => ({
    type: DeleteReviewNotesActionTypes.DELETE_REVIEW_NOTES,
    payload: {
        noteId, vaspId
    }
})

const updateReviewNoteApiResponseSuccess = (note) => ({
    type: UpdateReviewNotesActionTypes.API_RESPONSE_SUCCESS,
    payload: {
        note
    }
})


const createReviewNoteApiResponseSuccess = (note) => ({
    type: CreateReviewNoteActionTypes.API_RESPONSE_SUCCESS,
    payload: {
        note
    }
})



export { createReviewNoteApiResponseSuccess, updateReviewNoteApiResponseSuccess, deleteReviewNoteApiResponse, fetchReviewNotesApiResponse, fetchReviewNotesApiResponseError, fetchReviewNotesApiResponseSuccess }
