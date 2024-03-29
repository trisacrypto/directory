import produce from 'immer';

import { FetchVaspDetailsActionTypes, UpdateContactActionTypes } from './constants';

import {
  CreateReviewNoteActionTypes,
  DeleteContactActionTypes,
  DeleteReviewNotesActionTypes,
  FetchReviewNotesActionTypes,
  ReviewVaspActionTypes,
  UpdateBusinessInfosActionTypes,
  UpdateIvms101ActionTypes,
  UpdateReviewNotesActionTypes,
  UpdateTrisaImplementationDetailsActionTypes,
  UpdateTrixoActionTypes,
} from '.';

const INITIAL_STATE = {
  data: null,
  loading: false,
};

const vaspDetailsReducers = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case FetchVaspDetailsActionTypes.FETCH_VASP_DETAILS:
      return {
        ...state,
        loading: true,
      };
    case FetchVaspDetailsActionTypes.API_RESPONSE_SUCCESS:
      return {
        ...state,
        data: action.payload.data,
        loading: false,
      };
    case FetchVaspDetailsActionTypes.API_RESPONSE_ERROR:
      return {
        ...state,
        error: action.payload.error,
        loading: false,
      };
    case UpdateTrixoActionTypes.UPDATE_TRIXO:
      return {
        ...state,
        loading: true,
      };
    case UpdateTrixoActionTypes.API_RESPONSE_SUCCESS:
      return {
        ...state,
        loading: false,
        data: action.payload.data,
      };
    case UpdateTrixoActionTypes.API_RESPONSE_ERROR:
      return {
        ...state,
        loading: false,
        error: action.payload.error,
      };
    case UpdateBusinessInfosActionTypes.UPDATE_BUSINESS_INFOS:
    case UpdateTrisaImplementationDetailsActionTypes.UPDATE_TRISA_DETAILS:
    case UpdateIvms101ActionTypes.UPDATE_IVMS_101:
    case ReviewVaspActionTypes.REVIEW_VASP:
    case UpdateContactActionTypes.UPDATE_CONTACT:
    case DeleteContactActionTypes.DELETE_CONTACT:
      return {
        ...state,
        loading: true,
      };
    case UpdateBusinessInfosActionTypes.API_RESPONSE_SUCCESS:
    case UpdateTrisaImplementationDetailsActionTypes.API_RESPONSE_SUCCESS:
      return {
        ...state,
        loading: false,
        data: action.payload.data,
      };
    case UpdateBusinessInfosActionTypes.API_RESPONSE_ERROR:
      return {
        ...state,
        loading: false,
        error: action.payload.error,
      };
    case UpdateTrisaImplementationDetailsActionTypes.API_RESPONSE_ERROR:
      return {
        ...state,
        loading: false,
        trisaError: {
          message: action.payload.error.message,
          status: action.payload.error.errorStatus,
          statusText: action.payload.error.statusText,
        },
      };
    case UpdateTrisaImplementationDetailsActionTypes.CLEAR_ERROR_MESSAGE:
      return {
        ...state,
        loading: false,
        trisaError: null,
      };
    case UpdateIvms101ActionTypes.API_RESPONSE_SUCCESS:
      return {
        ...state,
        loading: false,
        data: action.payload.data,
      };
    case UpdateIvms101ActionTypes.API_RESPONSE_ERROR:
      return {
        ...state,
        loading: false,
        ivmsError: action.payload.error,
      };
    case UpdateIvms101ActionTypes.CLEAR_ERROR_MESSAGE:
      return {
        ...state,
        ivmsError: null,
      };
    case ReviewVaspActionTypes.API_RESPONSE_SUCCESS:
      return produce(state, (draft) => {
        draft.data.vasp.verification_status = action.payload.status;
      });
    case UpdateContactActionTypes.API_RESPONSE_ERROR:
      return {
        ...state,
        contactError: action.payload.error,
        loading: false,
      };
    case UpdateContactActionTypes.CLEAR_ERROR_MESSAGE:
      return {
        ...state,
        contactError: null,
        loading: false,
      };
    default:
      return state;
  }
};

const reviewNotesReducers = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case FetchReviewNotesActionTypes.FETCH_REVIEW_NOTES:
      return {
        ...state,
        loading: true,
      };
    case FetchReviewNotesActionTypes.API_RESPONSE_SUCCESS:
      return {
        ...state,
        data: action.payload.data,
        loading: false,
      };
    case FetchReviewNotesActionTypes.API_RESPONSE_ERROR:
      return {
        ...state,
        error: action.payload.error,
      };
    case DeleteReviewNotesActionTypes.DELETE_REVIEW_NOTES:
      const filteredData = state?.data?.filter((note) => note && note.id !== action.payload.noteId);
      return {
        ...state,
        data: filteredData,
      };
    case UpdateReviewNotesActionTypes.API_RESPONSE_SUCCESS:
      return produce(state, (draft) => {
        const idx = draft.data.findIndex((note) => note?.id === action.payload.note?.id);

        if (idx !== -1) {
          draft.data[idx] = action.payload.note;
        }
      });
    case CreateReviewNoteActionTypes.API_RESPONSE_SUCCESS:
      return produce(state, (draft) => {
        draft.data.unshift(action.payload.note);
      });
    default:
      return state;
  }
};

export { reviewNotesReducers, vaspDetailsReducers };
