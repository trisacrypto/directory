import { combineReducers } from 'redux';

import Layout from './layout/reducers';
import { vaspsReducers as Vasps, certificatesReducers as Certificates, summaryReducers as Summary, registrationsReviewsReducers as Reviews } from "./dashboard/reducers"
import { autocompleteReducers as Autocomplete } from './autocomplete';
import { vaspDetailsReducers as VaspDetails, reviewNotesReducers as ReviewNotes } from './vasp-details';

export default (combineReducers({
    Layout,
    Vasps,
    Certificates,
    Summary,
    Reviews,
    ReviewNotes,
    Autocomplete,
    VaspDetails
}));
