import React from "react";
import jwtDecode from "jwt-decode";
import { AUTH_SESSION_KEY } from "./constants";
import AuthContext from './auth-context';
import useSessionStorageState from "../../hooks/useSessionStorage";
import { APICore, setAuthorization, setCookie } from "../../helpers/api/apiCore";
import { getCookie } from "../../utils";

const api = new APICore()

const AuthProvider = ({ children }) => {
    const [token] = useSessionStorageState(AUTH_SESSION_KEY, '');
    const [authState, setAuthState] = React.useState(() => {
        if (token) {
            const tokenDecoded = jwtDecode(token.access_token);
            return {
                user: tokenDecoded,
                token: token
            }
        } else {
            return {}
        }
    });

    const setAuthInfo = (token) => {
        const tokenDecoded = jwtDecode(token.access_token)
        const csrfToken = getCookie('csrf_token')
        api.setLoggedInUser(token)
        setAuthorization(token.access_token)
        setCookie(csrfToken)
        setAuthState({
            user: tokenDecoded
        })


    }


    const value = {
        ...authState,
        setAuthInfo,
    }

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export default AuthProvider;