import dayjs from 'dayjs';
import jwtDecode from 'jwt-decode';
import React from 'react';

import { APICore, setAuthorization } from '@/helpers/api/apiCore';
import { postCredentials } from '@/helpers/api/auth';
import useSafeDispatch from '@/hooks/useSafeDispatch';

import AuthContext from './auth-context';
import { AUTH_SESSION_KEY } from './constants';

const api = new APICore();
const user = JSON.parse(sessionStorage.getItem(AUTH_SESSION_KEY));

const initialState = user
    ? {
          user: jwtDecode(user.access_token),
          isLoggedIn: true,
      }
    : {
          user: null,
          isLoggedIn: false,
      };

const reducer = (state = initialState, action) => {
    switch (action.type) {
        case 'LOGIN_SUCCESS':
            return {
                ...state,
                user: action.payload.user,
                isLoggedIn: true,
            };
        case 'LOGOUT':
            return {
                ...state,
                user: null,
                isLoggedIn: false,
            };
        default:
            break;
    }
};

const AuthProvider = ({ children }) => {
    const [state, _dispatch] = React.useReducer(reducer, initialState);
    const dispatch = useSafeDispatch(_dispatch);

    const isUserAuthenticated = () => {
        const { user } = state;
        if (!user) {
            return false;
        }

        return dayjs().diff(dayjs(user.exp)) > 0;
    };

    const login = async (user) => {
        try {
            const response = await postCredentials(user);
            const { access_token } = response.data;

            setAuthorization(access_token);
            api.setLoggedInUser(response.data);

            dispatch({
                type: 'LOGIN_SUCCESS',
                payload: { user: jwtDecode(access_token) },
            });
        } catch (error) {
            return Promise.reject(error);
        }
    };

    const logout = () => {
        window.google && window.google.accounts && window.google.accounts?.id.disableAutoSelect();
        dispatch({ type: 'LOGOUT' });
        api.deleteUserSession();
        api.setLoggedInUser(null);
    };

    const value = {
        ...state,
        dispatch,
        isUserAuthenticated,
        logout,
        login,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export default AuthProvider;
