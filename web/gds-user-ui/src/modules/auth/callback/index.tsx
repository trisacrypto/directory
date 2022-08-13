import React, { useEffect, useState } from 'react';

import { Box, useToast } from '@chakra-ui/react';
import useHashQuery from 'hooks/useHashQuery';
import { getAuth0User, userSelector } from 'modules/auth/login/user.slice';
import { clearCookies, getCookie } from 'utils/cookies';

import AlertMessage from 'components/ui/AlertMessage';
import { useNavigate } from 'react-router-dom';
import Loader from 'components/Loader';
import { t } from '@lingui/macro';
import { useSelector, useDispatch } from 'react-redux';
import * as Sentry from '@sentry/browser';
import ErrorMessage from 'components/ui/ErrorMessage';
const CallbackPage: React.FC = () => {
  const [isLoading, setIsloading] = useState(false);

  const query = useHashQuery();
  const accessToken = query.access_token;
  const callbackError = query.error;
  const { isFetching, isLoggedIn, isError, errorMessage } = useSelector(userSelector);

  const navigate = useNavigate();
  const dispatch = useDispatch();
  const toast = useToast();
  const getErrorMessage = () => {
    if (errorMessage?.error) {
      return errorMessage.error;
    } else {
      return ErrorMessage;
    }
  };
  useEffect(() => {
    dispatch(getAuth0User(accessToken));
  }, [accessToken]);

  useEffect(() => {
    if (callbackError) {
      toast({
        description: query.error_description,
        status: 'error',
        duration: 5000,
        isClosable: true,
        position: 'top-right'
      });
    }
    if (isError) {
      toast({
        description: getErrorMessage(),
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
  }, [isError, isLoggedIn, callbackError]);

  useEffect(() => {
    if (isFetching) {
      setIsloading(true);
    } else {
      setIsloading(false);
    }
  }, [isFetching]);

  return (
    <Box height={'100%'}>
      {isLoading && <Loader />}
      {isFetching && <Loader text="Loading Dashboard ..." />}
      {isError && (
        <AlertMessage
          title={callbackError || t`Token not valid`}
          message={query.error_description || getErrorMessage()}
          status="error"
        />
      )}
    </Box>
  );
};

export default CallbackPage;
