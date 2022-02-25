import React from 'react';
import { Redirect, useLocation } from 'react-router-dom';
import AccountLayout from './AccountLayout';
import SignWithGoogle from 'components/Auth/SignWithGoogle';
import config from "config";
import { APICore, setCookie } from 'helpers/api/apiCore';
import { getCookie } from 'utils';
import toast from 'react-hot-toast';
import { Alert, Col, Row } from 'react-bootstrap';
import useAuth from 'contexts/auth/use-auth';


const api = new APICore()

const Login = () => {

    const [csrfProtected, setCsrfProtected] = React.useState(false)
    const [loginError, setLoginError] = React.useState('');
    const isMounted = React.useRef(true)
    const { state } = useLocation()
    const { login } = useAuth()
    const [redirectOnLogin, setRedirectOnLogin] = React.useState(false)

    React.useEffect(() => {
        if (isMounted.current) {
            window.onGoogleLibraryLoad = () => {

                api.get('/authenticate').then(response => {
                    const csrfToken = getCookie('csrf_token')

                    setCookie(csrfToken);
                    setCsrfProtected(true)

                }).catch(error => {
                    toast.error(error)
                    console.error('[LOGIN] error:', error.message)
                })
            }
        }

        return () => { isMounted.current = false }
    }, [])

    const handleCredentialResponse = (response) => {
        if (response && response.credential) {
            const data = {
                credential
                    : response.credential
            }
            login(data).then(res => {
                setRedirectOnLogin(true)
            }).catch(error => {
                setLoginError("Address could not be authenticated (not a @trisa.io account).")
                console.error('[handleCredentialResponse]', error)
            })
        }

    }

    return (
        <>
            {redirectOnLogin ? <Redirect to={state ? state.from : '/'} /> : null}
            <AccountLayout >
                <Row>
                    <Col sm={12}>
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
                        {csrfProtected ? <SignWithGoogle clientId={config.GOOGLE_CLIENT_ID} text="signin_with" loginResponse={handleCredentialResponse} /> : <p className="text-center">loading...</p>}
                    </Col>
                </Row>
            </AccountLayout>
        </>
    );
};

export default Login;