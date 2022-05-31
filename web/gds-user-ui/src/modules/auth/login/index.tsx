import React, { useEffect, useState } from 'react';
import { Heading, Stack } from '@chakra-ui/react';
import Login from 'components/Section/Login';

import LandingLayout from 'layouts/LandingLayout';
import Head from 'components/Head/LandingHead';
import useCustomAuth0 from 'hooks/useCustomAuth0';
const StartPage: React.FC = () => {
  const [isLoading, setIsloading] = useState(false);
  const [error, setError] = useState('');
  const { auth0SignIn, auth0SignWithSocial, auth0Hash } = useCustomAuth0();

  const handleSocialAuth = (evt: any, type: any) => {
    evt.preventDefault();
    if (type === 'google') {
      auth0SignWithSocial('google-oauth2');
    }
  };
  const handleSignInWithEmail = async (data: any) => {
    setIsloading(true);
    try {
      const response: any = await auth0SignIn({
        username: data.username,
        password: data.password,
        responseType: 'token id_token',
        realm: 'Username-Password-Authentication'
      });
      if (response) {
        console.log('response ', response);
        setIsloading(false);
        if (response.emailVerified) {
          // to implement later
        }
      }
    } catch (err: any) {
      setIsloading(false);
      setError(err.code);
      // catch this error in sentry
      console.log('error', err);
    }
  };
  return (
    <LandingLayout>
      <Login
        handleSignWithSocial={handleSocialAuth}
        handleSignWithEmail={handleSignInWithEmail}
        isLoading={isLoading}
        isError={error}
      />
    </LandingLayout>
  );
};

export default StartPage;
