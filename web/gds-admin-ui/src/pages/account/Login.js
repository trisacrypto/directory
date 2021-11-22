import React from 'react';
import { Redirect, useLocation } from 'react-router-dom';
import AccountLayout from './AccountLayout';
import SignWithGoogle from 'components/Auth/SignWithGoogle';
import config from "config";
import { APICore, setCookie } from 'helpers/api/apiCore';
import { getCookie } from 'utils';
import toast from 'react-hot-toast';
import useAuth from 'contexts/auth/use-auth';
import { postCredentials } from 'helpers/api/auth';
import { Alert } from 'react-bootstrap';


const api = new APICore()

const Login = () => {

    const [csrfProtected, setCsrfProtected] = React.useState(false)
    const { setAuthInfo } = useAuth()
    const [redirectOnLogin, setRedirectOnLogin] = React.useState(
        false
    );
    const [loginError, setLogginError] = React.useState('');
    const isMounted = React.useRef(true)
    const { state } = useLocation()

    React.useEffect(() => {
        if (isMounted) {
            window.onload = () => {

                api.get('/authenticate').then(response => {
                    const csrfToken = getCookie('csrf_token')

                    setCookie(csrfToken);
                    setCsrfProtected(true)

                }).catch(error => {
                    toast.error(error)
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
            postCredentials(data).then(res => {
                setAuthInfo(res.data)
                setRedirectOnLogin(true)
            }).catch(error => {
                setLogginError(error)
                console.log('[Login] handleCredentialResponse', error)
            })
        }

    }

    return (
        <>
            {redirectOnLogin ? <Redirect to={state ? state.from : '/'}></Redirect> : null}
            <AccountLayout >
                {loginError && (
                    <Alert variant="danger" className="my-2">
                        {loginError}
                    </Alert>
                )}
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