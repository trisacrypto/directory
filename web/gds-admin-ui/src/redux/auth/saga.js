import { all, fork, put, takeEvery } from 'redux-saga/effects';

import { APICore } from '../../helpers/api/apiCore';
import { loginUserError, loginUserSuccess } from './actions';
import { AuthActionTypes } from './constants';
import jwtDecode from 'jwt-decode'

const api = new APICore();

/**
 * Login the user
 * @param {*} payload -token
 */
function* login({ payload }) {
    try {
        const user = jwtDecode(payload);
        api.setLoggedInUser(user);
        yield put(loginUserSuccess(AuthActionTypes.LOGIN_USER_SUCCESS, user));
    } catch (error) {
        yield put(loginUserError(AuthActionTypes.LOGIN_USER_ERROR, error));
        api.setLoggedInUser(null);
    }
}


export function* watchLoginUser() {
    yield takeEvery(AuthActionTypes.LOGIN_USER, login);
}



function* authSaga() {
    yield all([
        fork(watchLoginUser),
    ]);
}

export default authSaga;
