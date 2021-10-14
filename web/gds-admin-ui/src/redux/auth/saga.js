import { all, fork, put, takeEvery, call } from 'redux-saga/effects';

import { APICore, setAuthorization } from '../../helpers/api/apiCore';
import { loginUserError, loginUserSuccess } from './actions';
import { AuthActionTypes } from './constants';
import jwtDecode from 'jwt-decode'
import { postCredentials } from '../../helpers/api/auth';
import toast from 'react-hot-toast';

const api = new APICore();

/**
 * Login the user
 * @param {*} payload -token
 */
function* login({ payload }) {
    try {
        const response = yield call(postCredentials, payload)

        const data = response.data
        api.setLoggedInUser(data)
        setAuthorization(data.access_token)

        const decodedToken = jwtDecode(data.access_token)

        yield put(loginUserSuccess(AuthActionTypes.LOGIN_USER_SUCCESS, decodedToken));
    } catch (error) {
        console.log(error)
        toast.error(error)
        yield put(loginUserError(AuthActionTypes.LOGIN_USER_ERROR, error));
        api.setLoggedInUser(null);
        setAuthorization(null);
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
