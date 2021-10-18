// @flow
import { all } from 'redux-saga/effects';
import dashboardSaga, { vaspsSaga } from './dashboard/saga';
import layoutSaga from './layout/saga';

export default function* rootSaga(): any {
    yield all([layoutSaga(), dashboardSaga(), vaspsSaga()]);
}
