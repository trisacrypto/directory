import React, { useEffect, useState } from 'react';
import LandingLayout from 'layouts/LandingLayout';
import { Heading, Stack, Spinner } from '@chakra-ui/react';
import useHashQuery from 'hooks/useHashQuery';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import Cookies from 'universal-cookie';
import AlertMessage from 'components/ui/AlertMessage';
import { useNavigate } from 'react-router-dom';

const CallbackPage: React.FC = () => {
  const query = useHashQuery();
  const { auth0GetUser } = useCustomAuth0();
  const accessToken = query.access_token;
  const cookies = new Cookies();
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<any>('');
  useEffect(() => {
    (async () => {
      try {
        const getUserInfo: any = accessToken && (await auth0GetUser(accessToken));

        setIsLoading(false);
        if (getUserInfo && getUserInfo?.email_verified) {
          cookies.set('access_token', accessToken, { path: '/' });
          cookies.set('user_locale', getUserInfo?.locale, { path: '/' });
          console.log('getUserInfo', getUserInfo);
          navigate('/dashboard');
        } else {
          setError(
            'Your account has not been verified. Please check your email to verify your account.'
          );
        }
      } catch (e: any) {
        console.error('error', e);
      } finally {
        setIsLoading(false);
      }
    })();
  });

  return (
    <LandingLayout>
      {isLoading && <Spinner size={'2xl'} />}
      {error && <AlertMessage title="Token not valid" message={error} status="error" />}
    </LandingLayout>
  );
};

export default CallbackPage;
