import React, { useState } from 'react';
import { Heading, Stack } from '@chakra-ui/react';
import Register from 'components/Section/CreateAccount';

import LandingLayout from 'layouts/LandingLayout';
import useCustomAuth0 from 'hooks/useCustomAuth0';

const StartPage: React.FC = () => {
  const [isLoading, setIsloading] = useState(false);
  const [error, setError] = useState('');
  const [isPasswordError, setIsPasswordError] = useState(false);
  const [isUsernameError, setIsUsernameError] = useState(false);
  const { auth0SignUpWithEmail, auth0SignWithSocial } = useCustomAuth0();
  const handleSocialAuth = (evt: any, type: any) => {
    evt.preventDefault();
    if (type === 'google') {
      auth0SignWithSocial('google-oauth2');
    }
  };
  const handleSignUpWithEmail = async (data: any) => {
    console.log('datanfromform', data);
    setIsloading(true);
    try {
      const response: any = await auth0SignUpWithEmail({
        email: data.username,
        password: data.password,
        connection: 'Username-Password-Authentication'
      });
      if (response) {
        console.log('response', response);
        setIsloading(false);
        if (response.emailVerified) {
        }
      }
    } catch (err: any) {
      setIsloading(false);
      setError(err.code);
      if (err.code === 'invalid_password') {
        setIsPasswordError(true);
      }
      if (err.code === 'invalid_signup') {
        setIsUsernameError(true);
        console.log('username already exist');
      }
      // catch this error in sentry
      console.log('error', err);
    }
  };
  return (
    <LandingLayout>
      <Register
        handleSignUpWithEmail={handleSignUpWithEmail}
        handleSocialAuth={handleSocialAuth}
        isLoading={isLoading}
        isError={error}
        isPasswordError={isPasswordError}
        isUsernameError={isUsernameError}
      />
    </LandingLayout>
  );
};

export default StartPage;
