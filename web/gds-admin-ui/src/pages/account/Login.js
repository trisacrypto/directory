// @flow
import React from 'react';
import { Redirect } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';

import AccountLayout from './AccountLayout';
import SignWithGoogle from '../../components/Auth/SignWithGoogle';
import { loginUser } from '../../redux/auth/actions';
import config from "../../config";


const Login = (): React$Element<any> => {
    const dispatch = useDispatch();
    const { userIsLoggedIn, user } = useSelector((state) => ({
        userIsLoggedIn: state.Auth.userIsLoggedIn,
        user: state.Auth.user
    }))

    const handleCredentialResponse = (response) => {
        if (response.credential) {
            dispatch(loginUser(response.credential))
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
                <SignWithGoogle clientId={config.GOOGLE_CLIENT_ID} text="sign_in_with" loginResponse={handleCredentialResponse} />
            </AccountLayout>
        </>
    );
};

export default Login;
