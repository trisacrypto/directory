import { FetchCertificatesActionTypes, FetchPendingVaspsActionTypes, FetchRegistrationsReviewsActionTypes, FetchSummaryActionTypes, FetchVaspsActionTypes } from "./constants";


type Action = { type: string, payload: { data?: any, error?: string } };
type State = { data: {} | null, loading?: boolean };

const INITIAL_STATE = {
    data: null,
    loading: false
}

const SUMMARY_INITIAL_STATE = {
    data: null,
    loading: false
}

const certificatesReducers = (state: State = INITIAL_STATE, action: Action) => {

    switch (action.type) {
        case FetchCertificatesActionTypes.FETCH_CERTIFICATES:
            return {
                ...state,
                loading: true
            }
        case FetchCertificatesActionTypes.API_RESPONSE_SUCCESS:
            return {
                ...state,
                data: action.payload.data,
                loading: false
            }
        case FetchCertificatesActionTypes.API_RESPONSE_ERROR:
            return {
                ...state,
                error: action.payload.error
            }
        default:
            return state;
    }
}

const vaspsReducers = (state = INITIAL_STATE, action) => {

    switch (action.type) {
        case FetchVaspsActionTypes.FETCH_VASPS:
            return {
                ...state,
                loading: true
            }
        case FetchVaspsActionTypes.API_RESPONSE_SUCCESS:
            return {
                ...state,
                data: action.payload.data,
                loading: false
            }
        case FetchVaspsActionTypes.API_RESPONSE_ERROR:
            return {
                ...state,
                loading: false,
                error: action.payload.error
            }
        case FetchPendingVaspsActionTypes.FETCH_PENDING_VASPS:
            return {
                ...state,
                loading: true
            }
        default:
            return state;
    }
}


const summaryReducers = (state: State = SUMMARY_INITIAL_STATE, action: Action) => {

    switch (action.type) {
        case FetchSummaryActionTypes.FETCH_SUMMARY:
            return {
                ...state,
                loading: true
            }
        case FetchSummaryActionTypes.API_RESPONSE_SUCCESS:
            return {
                ...state,
                data: action.payload.data,
                loading: false
            }
        case FetchSummaryActionTypes.API_RESPONSE_ERROR:
            return {
                ...state,
                error: action.payload.error
            }
        default:
            return state;
    }
}



const registrationsReviewsReducers = (state: State = INITIAL_STATE, action: Action) => {

    switch (action.type) {
        case FetchRegistrationsReviewsActionTypes.FETCH_REGISTRATIONS_REVIEWS:
            return {
                ...state,
                loading: true
            }
        case FetchRegistrationsReviewsActionTypes.API_RESPONSE_SUCCESS:
            return {
                ...state,
                data: action.payload.data,
                loading: false
            }
        case FetchRegistrationsReviewsActionTypes.API_RESPONSE_ERROR:
            return {
                ...state,
                error: action.payload.error
            }
        default:
            return state;
    }
}

export { certificatesReducers, vaspsReducers, summaryReducers, registrationsReviewsReducers };