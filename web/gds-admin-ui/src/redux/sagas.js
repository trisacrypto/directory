import { all } from 'redux-saga/effects';

import autocompletesSaga from './autocomplete/saga';
import dashboardSaga, { vaspsSaga } from './dashboard/saga';
import layoutSaga from './layout/saga';
import {
  deleteContactSaga,
  deleteReviewNoteSaga,
  reviewNotesSaga,
  updateBusinessInfosSaga,
  updateContactSaga,
  updateIvmsSaga,
  updateTrisaDetailsSaga,
  updateTrixoSaga,
  vaspDetailsSaga,
} from './vasp-details';

export default function* rootSaga() {
  yield all([
    layoutSaga(),
    dashboardSaga(),
    vaspsSaga(),
    reviewNotesSaga(),
    deleteReviewNoteSaga(),
    autocompletesSaga(),
    vaspDetailsSaga(),
    updateTrixoSaga(),
    updateBusinessInfosSaga(),
    updateTrisaDetailsSaga(),
    updateIvmsSaga(),
    updateContactSaga(),
    deleteContactSaga(),
  ]);
}
