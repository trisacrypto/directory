import { combineReducers } from 'redux';

import { autocompleteReducers as Autocomplete } from './autocomplete';
import {
  certificatesReducers as Certificates,
  registrationsReviewsReducers as Reviews,
  summaryReducers as Summary,
  vaspsReducers as Vasps,
} from './dashboard/reducers';
import Layout from './layout/reducers';
import {
  reviewNotesReducers as ReviewNotes,
  vaspDetailsReducers as VaspDetails,
} from './vasp-details';

export default combineReducers({
  Layout,
  Vasps,
  Certificates,
  Summary,
  Reviews,
  ReviewNotes,
  Autocomplete,
  VaspDetails,
});
