// @flow
import { combineReducers } from 'redux';

import Layout from './layout/reducers';
import { vaspsReducers as Vasps, certificatesReducers as Certificates, summaryReducers as Summary } from "./dashboard/reducers"

export default (combineReducers({
    Layout,
    Vasps,
    Certificates,
    Summary
}): any);
