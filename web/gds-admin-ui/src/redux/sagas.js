// @flow
import { all } from 'redux-saga/effects';
import dashboardSaga from './dashboard/saga';
import layoutSaga from './layout/saga';
import authSaga from './auth/saga';

export default function* rootSaga(): any {
    yield all([layoutSaga(), dashboardSaga(), authSaga()]);
}
