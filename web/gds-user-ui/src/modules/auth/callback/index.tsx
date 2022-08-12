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
const CallbackPage: React.FC = () => {
  const [isLoading, setIsloading] = useState(false);
  const query = useHashQuery();
  const accessToken = query.access_token;
  const { isFetching, isLoggedIn, isError, errorMessage } = useSelector(userSelector);
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const toast = useToast();
  useEffect(() => {
    clearCookies();
    dispatch(getAuth0User(accessToken));
  }, [accessToken, dispatch]);

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
      {isError && <AlertMessage title={t`Token not valid`} message={errorMessage} status="error" />}
    </Box>
  );
};

export default CallbackPage;
