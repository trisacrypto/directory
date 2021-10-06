// @flow
import React from 'react';
import { Redirect } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';

import AccountLayout from './AccountLayout';
import SignWithGoogle from '../../components/Auth/SignWithGoogle';
import { loginUser } from '../../redux/auth/actions';
import config from "../../config";
import { APICore, setCookie } from '../../helpers/api/apiCore';
import { getCookie } from '../../utils';

const api = new APICore()

const Login = (): React$Element<any> => {
    const dispatch = useDispatch();
    const { userIsLoggedIn, user } = useSelector((state) => ({
        userIsLoggedIn: state.Auth.userIsLoggedIn,
        user: state.Auth.user
    }))
    const [csrfProtected, setCsrfProtected] = React.useState(false)
    const isMounted = React.useRef(true)


    React.useEffect(() => {
        if (isMounted) {
            window.onload = () => {

                api.get('/authenticate').then(response => {
                    const csrfToken = getCookie('csrf_token')

                    setCookie(csrfToken);
                    setCsrfProtected(true)

                }).catch(error => {
                    console.log('[Login] error:', error.message)
                })
            }
        }

        return () => { isMounted.current = false }

    }, [])

    const handleCredentialResponse = (response) => {
        if (response.credential) {
            const data = {
                credential
                    : response.credential
            }
            dispatch(loginUser(data))
        }

    }

    return (
        <>
            {userIsLoggedIn || user ? <Redirect to="/" /> : null}
            <AccountLayout >
                <div className="text-center w-75 m-auto">
                    <h4 className="text-dark-50 text-center mt-0 fw-bold">Sign In</h4>
                    <p className="text-muted mb-4">
                        Please use your @trisa.io Google account to access the GDS Admin.
                    </p>
                </div>
                {csrfProtected ? <SignWithGoogle clientId={config.GOOGLE_CLIENT_ID} text="sign_in_with" loginResponse={handleCredentialResponse} /> : <p className="text-center">loading...</p>}
            </AccountLayout>
        </>
    );
};

export default Login;