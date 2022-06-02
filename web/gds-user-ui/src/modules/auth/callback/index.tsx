import React, { useEffect, useState } from 'react';

import { Heading, Stack, Spinner, Flex, Box } from '@chakra-ui/react';
import useHashQuery from 'hooks/useHashQuery';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import { getCookie, setCookie } from 'utils/cookies';
import AlertMessage from 'components/ui/AlertMessage';
import { useNavigate } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
import { t } from '@lingui/macro';

const CallbackPage: React.FC = () => {
  const query = useHashQuery();
  const { auth0GetUser } = useCustomAuth0();
  const { loginUser } = useAuth();
  const accessToken = query.access_token;
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<any>('');
  useEffect(() => {
    (async () => {
      try {
        const getUserInfo: any = accessToken && (await auth0GetUser(accessToken));

        setIsLoading(false);
        if (getUserInfo && getUserInfo?.email_verified) {
          setCookie('access_token', accessToken);
          setCookie('user_locale', getUserInfo?.locale);
          const userInfo: TUser = {
            isLoggedIn: true,
            user: {
              name: getUserInfo?.name,
              pictureUrl: getUserInfo?.picture,
              email: getUserInfo?.email
            }
          };
          loginUser(userInfo);
          navigate('/dashboard/overview');
        } else {
          setError(
            t`Your account has not been verified. Please check your email to verify your account.`
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
    <Box height={'100%'}>
      {isLoading && (
        <Box
          textAlign={'center'}
          justifyItems="center"
          alignItems={'center'}
          justifyContent="center">
          <Spinner size={'xl'} />
        </Box>
      )}
      {error && <AlertMessage title={t`Token not valid`} message={error} status="error" />}
    </Box>
  );
};

export default CallbackPage;
