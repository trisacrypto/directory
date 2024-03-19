import React, { useEffect, useState } from 'react';
import Login from 'components/Section/Login';
import useAuth from 'hooks/useAuth';
import LandingLayout from 'layouts/LandingLayout';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import useSearchParams from 'hooks/useQueryParams';
import * as Sentry from '@sentry/browser';
import useCustomToast from 'hooks/useCustomToast';
import useCertificateStepper from 'hooks/useCertificateStepper';
const StartPage: React.FC = () => {
  const [isLoading, setIsloading] = useState(false);
  const [error, setError] = useState('');
  const { auth0SignIn, auth0SignWithSocial } = useCustomAuth0();
  const { clearStepperState } = useCertificateStepper();
  const { loginUser } = useAuth();
  const { q, error_description, orgid } = useSearchParams();
  const toast = useCustomToast();
  useEffect(() => {
    // rend tost if q is not empty
    let message = '';
    if (q) {
      if (q === 'unauthorized') {
        message =
          'Your account does not have permission to access the Administrator interface. Contact the administrator of your organization for assistance';
      }
      if (q === 'token_expired' || q === 'invalid_token') {
        message =
          'Your session has expired. Please sign in again to continue using the Administrator interface';
      }
    }
    if (error_description) {
      message = error_description;
    }
    if (message) {
      toast({
        description: message
      });
    }
    if (orgid) {
      localStorage.setItem('orgId', orgid);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [q, error_description, orgid]);

  // clean cookies

  useEffect(() => {
    if (isLoading) {
      clearStepperState();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isLoading]);

  const handleSocialAuth = (evt: any, type: any) => {
    evt.preventDefault();

    if (type === 'google') {
      auth0SignWithSocial('google-oauth2');
    }
    if (type === 'github') {
      auth0SignWithSocial('github');
    }

    if (type === 'microsoft') {
      auth0SignWithSocial('windowslive');
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
          // get user info
          loginUser(response);
        } else {
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
      Sentry.captureException(err?.original);
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
