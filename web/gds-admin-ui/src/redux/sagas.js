// @flow
import { all } from 'redux-saga/effects';
import dashboardSaga from './dashboard/saga';

import layoutSaga from './layout/saga';

export default function* rootSaga(): any {
    yield all([layoutSaga(), dashboardSaga()]);
}
