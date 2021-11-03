import { all } from 'redux-saga/effects';
import dashboardSaga, { vaspsSaga } from './dashboard/saga';
import layoutSaga from './layout/saga';
import { createReviewNoteSaga, reviewNotesSaga } from './review-notes';

export default function* rootSaga() {
    yield all([layoutSaga(), dashboardSaga(), vaspsSaga(), reviewNotesSaga()]);
}
