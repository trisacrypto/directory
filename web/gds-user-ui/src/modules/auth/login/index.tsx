import React, { useEffect, useState } from 'react';
import { Heading, Stack, useToast } from '@chakra-ui/react';
import Login from 'components/Section/Login';
import * as Sentry from '@sentry/browser';
import LandingLayout from 'layouts/LandingLayout';
import Head from 'components/Head/LandingHead';
import useCustomAuth0 from 'hooks/useCustomAuth0';
const StartPage: React.FC = () => {
  const [isLoading, setIsloading] = useState(false);
  const [error, setError] = useState('');
  const { auth0SignIn, auth0SignWithSocial, auth0Hash } = useCustomAuth0();
  const toast = useToast();
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
        setIsloading(false);
        if (response.emailVerified) {
          // to implement later
        }
      }
    } catch (err: any) {
      setIsloading(false);
      if (err.code === 'access_denied') {
        toast({
          description: 'Invalid username or password',
          status: 'error',
          duration: 5000,
          position: 'top',
          isClosable: true
        });

        setError('Invalid username or password');
      }

      // catch this error in sentry
      Sentry.captureException(err);
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
