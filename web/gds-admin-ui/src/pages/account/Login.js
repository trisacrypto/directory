// @flow
import React from 'react';
import { Button } from 'react-bootstrap';
import { Link, useHistory } from 'react-router-dom';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';

import { VerticalForm, FormInput } from '../../components/';

import AccountLayout from './AccountLayout';


const Login = (): React$Element<any> => {
    const history = useHistory();

    const onSubmit = (formData) => {
        console.log("[Login] Form Date", formData)
        history.push("/dashboard")
    };

    const schemaResolver = yupResolver(
        yup.object().shape({
            username: yup.string().required('Please enter Username'),
            password: yup.string().required('Please enter Password')
        })
    );

    return (
        <>

            <AccountLayout >
                <div className="text-center w-75 m-auto">
                    <h4 className="text-dark-50 text-center mt-0 fw-bold">Sign In</h4>
                    <p className="text-muted mb-4">
                        Please use your @trisa.io Google account to access the GDS Admin.
                    </p>
                </div>

                <VerticalForm
                    onSubmit={onSubmit}
                    resolver={schemaResolver}
                    defaultValues={{ username: '', password: '' }}>
                    <FormInput
                        label='Username'
                        type="text"
                        name="username"
                        placeholder='Enter your Username'
                        containerClass={'mb-3'}
                    />
                    <FormInput
                        label='Password'
                        type="password"
                        name="password"
                        placeholder='Enter your password'
                        containerClass={'mb-3'}>
                        <Link to="/account/forget-password" className="text-muted float-end">
                            <small>Forgot your password?</small>
                        </Link>
                    </FormInput>

                    <div className="mb-3 mb-0 text-center">
                        <Button variant="primary" type="submit">
                            Log In
                        </Button>
                    </div>
                </VerticalForm>
            </AccountLayout>
        </>
    );
};

export default Login;
