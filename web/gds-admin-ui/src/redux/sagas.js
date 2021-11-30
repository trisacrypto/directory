import { all } from 'redux-saga/effects';
import autocompletesSaga from './autocomplete/saga';
import dashboardSaga, { vaspsSaga } from './dashboard/saga';
import layoutSaga from './layout/saga';
import { deleteReviewNoteSaga, reviewNotesSaga } from './review-notes';

export default function* rootSaga() {
    yield all([layoutSaga(), dashboardSaga(), vaspsSaga(), reviewNotesSaga(), deleteReviewNoteSaga(), autocompletesSaga()]);
}
