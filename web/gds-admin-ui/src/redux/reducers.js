// @flow
import { combineReducers } from 'redux';

import Layout from './layout/reducers';
import { vaspsReducers as Vasps, certificatesReducers as Certificates, summaryReducers as Summary, registrationsReviewsReducers as Reviews } from "./dashboard/reducers"
import Auth from './auth/reducers'

export default (combineReducers({
    Layout,
    Vasps,
    Certificates,
    Summary,
    Auth,
    Reviews
}): any);
