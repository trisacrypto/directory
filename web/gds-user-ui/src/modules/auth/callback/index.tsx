import React, { useEffect, useState } from 'react';

import { Heading, Stack, Spinner, Flex, Box, useToast } from '@chakra-ui/react';
import useHashQuery from 'hooks/useHashQuery';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import { getCookie, setCookie } from 'utils/cookies';
import AlertMessage from 'components/ui/AlertMessage';
import { useNavigate } from 'react-router-dom';
import useAuth from 'hooks/useAuth';
import { t } from '@lingui/macro';

import { logUserInBff } from 'modules/auth/login/auth.service';
const CallbackPage: React.FC = () => {
  const query = useHashQuery();
  const { auth0GetUser, auth0Hash } = useCustomAuth0();
  const { loginUser } = useAuth();
  const accessToken = query.access_token;
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<any>('');
  const toast = useToast();
  useEffect(() => {
    (async () => {
      try {
        const getUserInfo: any = accessToken && (await auth0Hash());
        console.log('[getUserInfo]', getUserInfo);
        setIsLoading(false);
        if (getUserInfo && getUserInfo?.email_verified) {
          setCookie('access_token', accessToken);
          setCookie('user_locale', getUserInfo?.locale);
          const getUser = await logUserInBff();

          // if (getUser.status === 204) {
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
          // }
          // log this error to sentry
        } else {
          setError(
            t`Your account has not been verified. Please check your email to verify your account.`
          );
        }
      } catch (e: any) {
        toast({
          description: e.response?.data?.message || e.message,
          status: 'error',
          duration: 5000,
          isClosable: true,
          position: 'top-right'
        });
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
