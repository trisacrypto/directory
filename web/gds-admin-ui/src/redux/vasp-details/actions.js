import { CreateReviewNoteActionTypes, DeleteReviewNotesActionTypes, FetchReviewNotesActionTypes, ReviewVaspActionTypes, UpdateBusinessInfosActionTypes, UpdateIvms101ActionTypes, UpdateReviewNotesActionTypes, UpdateTrisaImplementationDetailsActionTypes, UpdateTrixoActionTypes } from ".";
import { DeleteContactActionTypes, FetchVaspDetailsActionTypes, UpdateContactActionTypes } from "./constants";

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

const reviewVaspApiResponse = () => ({
    type: ReviewVaspActionTypes.REVIEW_VASP,
    payload: {}
})

const reviewVaspApiResponseSuccess = (status) => ({
    type: ReviewVaspActionTypes.API_RESPONSE_SUCCESS,
    payload: { status }
})


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

export const updateBusinessInfosResponse = (id, businessInfos, setIsOpen) => ({
    type: UpdateBusinessInfosActionTypes.UPDATE_BUSINESS_INFOS,
    payload: { id, businessInfos, setIsOpen }
})

export const updateBusinessInfosResponseSuccess = (data) => ({
    type: UpdateBusinessInfosActionTypes.API_RESPONSE_SUCCESS,
    payload: { data }
})

export const updateBusinessInfosResponseError = (error) => ({
    type: UpdateBusinessInfosActionTypes.API_RESPONSE_ERROR,
    payload: { error }
})

export const updateTrisaImplementationDetailsResponse = (id, trisa, setIsOpen) => ({
    type: UpdateTrisaImplementationDetailsActionTypes.UPDATE_TRISA_DETAILS,
    payload: { id, trisa, setIsOpen }
})

export const updateTrisaImplementationDetailsResponseSuccess = (data) => ({
    type: UpdateTrisaImplementationDetailsActionTypes.API_RESPONSE_SUCCESS,
    payload: { data }
})

export const updateTrisaImplementationDetailsResponseError = (error) => ({
    type: UpdateTrisaImplementationDetailsActionTypes.API_RESPONSE_ERROR,
    payload: { error }
})

export const clearTrisaImplementationDetailsErrorMessage = () => ({
    type: UpdateTrisaImplementationDetailsActionTypes.CLEAR_ERROR_MESSAGE,
    payload: {}
})

export const updateIvms101Response = (id, ivms, setIsOpen) => ({
    type: UpdateIvms101ActionTypes.UPDATE_IVMS_101,
    payload: { id, ivms, setIsOpen }
})

export const updateIvms101ResponseSuccess = (data) => ({
    type: UpdateIvms101ActionTypes.API_RESPONSE_SUCCESS,
    payload: { data }
})

export const updateIvms101ResponseError = (error) => ({
    type: UpdateIvms101ActionTypes.API_RESPONSE_ERROR,
    payload: { error }
})

export const clearIvms101ErrorMessage = () => ({
    type: UpdateIvms101ActionTypes.CLEAR_ERROR_MESSAGE,
    payload: {}
})

export const updateContact = ({ vaspId, contactType, data, setIsOpen }) => ({
    type: UpdateContactActionTypes.UPDATE_CONTACT,
    payload: { vaspId, contactType, data, setIsOpen }
})

export const updateContactResponseError = (error) => ({
    type: UpdateContactActionTypes.API_RESPONSE_ERROR,
    payload: { error }
})

export const deleteContactResponse = (vaspId, contactType, setIsOpen) => ({
    type: DeleteContactActionTypes.DELETE_CONTACT,
    payload: { vaspId, contactType, setIsOpen }
})

export const deleteContactResponseSuccess = (data) => ({
    type: DeleteContactActionTypes.DELETE_CONTACT,
    payload: { data }
})

export const deleteContactResponseError = (error) => ({
    type: DeleteContactActionTypes.API_RESPONSE_ERROR,
    payload: { error }
})



export {
    reviewVaspApiResponse,
    reviewVaspApiResponseSuccess,
    updateTrixoResponse,
    updateTrixoResponseSuccess,
    updateTrixoResponseError,
    fetchVaspDetailsApiResponse,
    fetchVaspDetailsApiResponseError,
    fetchVaspDetailsApiResponseSuccess,
    createReviewNoteApiResponseSuccess,
    updateReviewNoteApiResponseSuccess,
    deleteReviewNoteApiResponse,
    fetchReviewNotesApiResponse,
    fetchReviewNotesApiResponseError,
    fetchReviewNotesApiResponseSuccess
}
