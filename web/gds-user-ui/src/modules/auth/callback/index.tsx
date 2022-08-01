import React, { useEffect } from 'react';

import { Spinner, Box, useToast } from '@chakra-ui/react';
import useHashQuery from 'hooks/useHashQuery';
import { getAuth0User, userSelector } from 'modules/auth/login/user.slice';
import AlertMessage from 'components/ui/AlertMessage';
import { useNavigate } from 'react-router-dom';
import Loader from 'components/Loader';
import { t } from '@lingui/macro';
import { useDispatch, useSelector } from 'react-redux';

const CallbackPage: React.FC = () => {
  const query = useHashQuery();
  const accessToken = query.access_token;
  const { isFetching, isLoggedIn, isError, errorMessage } = useSelector(userSelector);
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const toast = useToast();
  useEffect(() => {
    dispatch(getAuth0User(accessToken));

    // (async () => {
    //   try {
    //     const getUserInfo: any = accessToken && (await auth0Hash());
    //     console.log('[getUserInfo]', getUserInfo);
    //     setIsLoading(false);
    //     if (getUserInfo && getUserInfo?.idTokenPayload.email_verified) {
    //       setCookie('access_token', accessToken);
    //       setCookie('user_locale', getUserInfo?.locale);
    //       const getUser = await logUserInBff();

    //       // if (getUser.status === 204) {
    //       const userInfo: TUser = {
    //         isLoggedIn: true,
    //         user: {
    //           name: getUserInfo?.name,
    //           pictureUrl: getUserInfo?.picture,
    //           email: getUserInfo?.email
    //         }
    //       };
    //       console.log('[login dispatch] second');
    //       loginUser(userInfo);
    //       navigate('/dashboard/overview');
    //       // }
    //       // log this error to sentry
    //     } else {
    //       setError(
    //         t`Your account has not been verified. Please check your email to verify your account.`
    //       );
    //     }
    //   } catch (e: any) {
    //     toast({
    //       description: e.response?.data?.message || e.message,
    //       status: 'error',
    //       duration: 5000,
    //       isClosable: true,
    //       position: 'top-right'
    //     });
    //   } finally {
    //     setIsLoading(false);
    //   }
    // })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [accessToken]);

  useEffect(() => {
    if (isError) {
      toast({
        description: errorMessage,
        status: 'error',
        duration: 5000,
        isClosable: true,
        position: 'top-right'
      });
    }
    if (isLoggedIn) {
      navigate('/dashboard/overview');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isError, isLoggedIn]);

  return (
    <Box height={'100%'}>
      {isFetching && <Loader text="Loading Dashboard ..." />}
      {isError && <AlertMessage title={t`Token not valid`} message={errorMessage} status="error" />}
    </Box>
  );
};

export default CallbackPage;
