// @flow
import { AuthActionTypes } from './constants';

import { APICore } from '../../helpers/api/apiCore';
import jwtDecode from 'jwt-decode';

const api = new APICore();

const user = api.getLoggedInUser()
const decodedUser = user ? jwtDecode(user.access_token) : ''


const INIT_STATE = {
    user: decodedUser,
    loading: false,
};


const Auth = (state = INIT_STATE, action) => {

    switch (action.type) {
        case AuthActionTypes.LOGIN_USER:
            return {
                ...state,
                loading: true
            }
        case AuthActionTypes.LOGIN_USER_SUCCESS:
            return {
                user: action.payload,
                userIsloggedIn: true,
                loading: false
            }
        case AuthActionTypes.LOGIN_USER_ERROR:
            return {
                error: action.payload.error,
                userIsloggedIn: false,
                loading: false
            }
        default:
            return { ...state };
    }
};

export default Auth;
