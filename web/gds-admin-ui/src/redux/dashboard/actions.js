import { FetchCertificatesActionTypes, FetchVaspsActionTypes, FetchSummaryActionTypes, FetchRegistrationsReviewsActionTypes } from "./constants";

type Action = { type: string, payload: {} | string };

const fetchSummaryApiResponseSuccess = (actionType: string, data: any): Action => ({
    type: FetchSummaryActionTypes.API_RESPONSE_SUCCESS,
    payload: { actionType, data },
});

const fetchSummaryApiResponseError = (actionType: string, data: any): Action => ({
    type: FetchSummaryActionTypes.API_RESPONSE_ERROR,
    payload: { actionType, data },
});

const fetchSummary = (): Action => ({
    type: FetchSummaryActionTypes.FETCH_SUMMARY,
    payload: {},
});


const fetchCertificateApiResponseSuccess = (actionType: string, data: any): Action => ({
    type: FetchCertificatesActionTypes.API_RESPONSE_SUCCESS,
    payload: { actionType, data },
});

const fetchCertificateApiResponseError = (actionType: string, data: any): Action => ({
    type: FetchCertificatesActionTypes.API_RESPONSE_ERROR,
    payload: { actionType, data },
});

const fetchCertificates = (): Action => ({
    type: FetchCertificatesActionTypes.FETCH_CERTIFICATES,
    payload: {},
});


const fetchVaspsApiResponseSuccess = (actionType: string, data: any): Action => ({
    type: FetchVaspsActionTypes.API_RESPONSE_SUCCESS,
    payload: { actionType, data },
});

const fetchVaspsApiResponseError = (actionType: string, data: any): Action => ({
    type: FetchVaspsActionTypes.API_RESPONSE_ERROR,
    payload: { actionType, data },
});

const fetchVasps = (): Action => ({
    type: FetchVaspsActionTypes.FETCH_VASPS,
    payload: {},
});

const fecthRegistrationsReviews = (): Action => ({
    type: FetchRegistrationsReviewsActionTypes.FETCH_REGISTRATIONS_REVIEWS,
    payload: {}
})

const fetchRegistrationsReviewsSuccess = (actionType: string, data: any): Action => ({
    type: FetchRegistrationsReviewsActionTypes.API_RESPONSE_SUCCESS,
    payload: {
        actionType, data
    }
})

const fetchRegistrationsReviewsError = (actionType: string, error: any): Action => ({
    type: FetchRegistrationsReviewsActionTypes.API_RESPONSE_ERROR,
    payload: {
        actionType, error
    }
})

export {
    fetchVasps, fetchVaspsApiResponseError, fetchVaspsApiResponseSuccess, fetchCertificates, fetchCertificateApiResponseError, fetchCertificateApiResponseSuccess, fetchSummaryApiResponseError, fetchSummaryApiResponseSuccess, fetchSummary, fetchRegistrationsReviewsError, fetchRegistrationsReviewsSuccess, fecthRegistrationsReviews
}