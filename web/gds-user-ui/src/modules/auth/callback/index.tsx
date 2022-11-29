import React, { useEffect, useState } from 'react';
import { Box } from '@chakra-ui/react';
import useHashQuery from 'hooks/useHashQuery';
import { getAuth0User, userSelector } from 'modules/auth/login/user.slice';
import { useNavigate } from 'react-router-dom';
import Loader from 'components/Loader';
import { useSelector, useDispatch } from 'react-redux';
const CallbackPage: React.FC = () => {
  const [isLoading, setIsloading] = useState(true);

  const query = useHashQuery();
  const { access_token: accessToken, error: callbackError, error_description } = query as any;
  const { isFetching, isLoggedIn, isError, errorMessage } = useSelector(userSelector);

  const navigate = useNavigate();
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(getAuth0User({ hasToken: accessToken }));
  }, [accessToken, dispatch]);

  useEffect(() => {
    if (callbackError || isError) {
      navigate(`/auth/login?error_description=${error_description || errorMessage}`);
    }
    if (isLoggedIn) {
      navigate('/dashboard/overview');
    }
  }, [isError, isLoggedIn, callbackError, errorMessage, navigate, error_description]);

  useEffect(() => {
    if (!isFetching) {
      setIsloading(true);
    } else {
      setIsloading(false);
    }
  }, [isFetching]);

  return (
    <Box height={'100%'}>
      {isLoading && <Loader />}
      {isFetching && <Loader text="Loading Dashboard ..." />}
    </Box>
  );
};

export default CallbackPage;
