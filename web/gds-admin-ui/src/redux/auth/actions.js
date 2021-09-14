import { AuthActionTypes } from './constants';

export const loginUser = (data) => ({
    type: AuthActionTypes.LOGIN_USER,
    payload: data
});

export const loginUserSuccess = (actionType, data) => {
    return {
        type: AuthActionTypes.LOGIN_USER_SUCCESS,
        payload: data,
    }
};


export const loginUserError = (error) => ({
    type: AuthActionTypes.LOGIN_USER_ERROR,
    payload: error,
});
