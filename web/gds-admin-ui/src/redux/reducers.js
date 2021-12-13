import { combineReducers } from 'redux';

import Layout from './layout/reducers';
import { vaspsReducers as Vasps, certificatesReducers as Certificates, summaryReducers as Summary, registrationsReviewsReducers as Reviews } from "./dashboard/reducers"
import { reviewNotesReducers as ReviewNotes } from './review-notes';
import { autocompleteReducers as Autocomplete } from './autocomplete';

export default (combineReducers({
    Layout,
    Vasps,
    Certificates,
    Summary,
    Reviews,
    ReviewNotes,
    Autocomplete
}));
