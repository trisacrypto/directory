// @flow
import { AuthActionTypes } from './constants';


const INIT_STATE = {
    user: null,
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
            console.log('SUCCESS', action)
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
