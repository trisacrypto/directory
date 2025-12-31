import React from 'react';
import { Alert, Col, Row } from 'react-bootstrap';
import { Redirect, useLocation } from 'react-router-dom';

import SignWithGoogle from '@/components/Auth/SignWithGoogle';
import config from '@/config';
import useAuth from '@/contexts/auth/use-auth';

import AccountLayout from './AccountLayout';
import { useGetAuthenticate } from '../../services';
import useScript from '@/hooks/useScript';
import { getCookie } from '@/utils';
import { captureException } from '@sentry/react';

const Login = () => {
    const [loginError, setLoginError] = React.useState('');
    const { state } = useLocation();
    const { login } = useAuth();
    const [redirectOnLogin, setRedirectOnLogin] = React.useState(false);
    const loadScriptStatus = useScript('https://accounts.google.com/gsi/client');
    // eslint-disable-next-line no-unused-vars
    const { data } = useGetAuthenticate({
        loadScriptStatus: loadScriptStatus,
    });

    const csrfProtected = !!getCookie('csrf_token');

    const handleCredentialResponse = (response) => {
        if (response && response.credential) {
            const data = {
                credential: response.credential,
            };
            login(data)
                .then((res) => {
                    setRedirectOnLogin(true);
                })
                .catch((error) => {
                    setLoginError('Address could not be authenticated (not a @travelrule.io account).');
                    captureException(error);
                });
        }
    };

    return (
        <>
            {redirectOnLogin ? <Redirect to={state ? state.from : '/'} /> : null}
            <AccountLayout>
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
                                Please use your @travelrule.io Google account to access the GDS Admin.
                            </p>
                        </div>
                        {csrfProtected && (
                            <SignWithGoogle
                                clientId={config.GOOGLE_CLIENT_ID}
                                text="signin_with"
                                loginResponse={handleCredentialResponse}
                            />
                        )}
                    </Col>
                </Row>
            </AccountLayout>
        </>
    );
};

export default Login;
