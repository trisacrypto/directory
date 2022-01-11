import { CreateReviewNoteActionTypes, DeleteReviewNotesActionTypes, FetchReviewNotesActionTypes, UpdateReviewNotesActionTypes, UpdateTrixoActionTypes } from ".";
import { FetchVaspDetailsActionTypes } from "./constants";

const fetchVaspDetailsApiResponse = (vaspId, history) => ({
    type: FetchVaspDetailsActionTypes.FETCH_VASP_DETAILS,
    payload: { id: vaspId, history },
});

const fetchVaspDetailsApiResponseSuccess = (data) => ({
    type: FetchVaspDetailsActionTypes.API_RESPONSE_SUCCESS,
    payload: { data },
});

const fetchVaspDetailsApiResponseError = (error) => ({
    type: FetchVaspDetailsActionTypes.API_RESPONSE_ERROR,
    payload: { error },
});


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

const updateTrixoResponse = (id, trixo, setIsOpen) => ({
    type: UpdateTrixoActionTypes.UPDATE_TRIXO,
    payload: { id, trixo, setIsOpen }
})

const updateTrixoResponseSuccess = (data) => ({
    type: UpdateTrixoActionTypes.API_RESPONSE_SUCCESS,
    payload: { data }
})

const updateTrixoResponseError = (error) => ({
    type: UpdateTrixoActionTypes.API_RESPONSE_ERROR,
    payload: { error }
})


export { updateTrixoResponse, updateTrixoResponseSuccess, updateTrixoResponseError, fetchVaspDetailsApiResponse, fetchVaspDetailsApiResponseError, fetchVaspDetailsApiResponseSuccess, createReviewNoteApiResponseSuccess, updateReviewNoteApiResponseSuccess, deleteReviewNoteApiResponse, fetchReviewNotesApiResponse, fetchReviewNotesApiResponseError, fetchReviewNotesApiResponseSuccess }