import React, { useEffect, useState } from 'react';
import { removeCookie } from 'utils/cookies';
import { Stack, Flex } from '@chakra-ui/react';
import Loader from 'components/Loader';
import { useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { logout } from 'modules/auth/login/user.slice';
import { setDefaultMemberNetwork } from 'modules/dashboard/member/member.slice';
import Store from 'application/store';
const Logout: React.FC = () => {
  const dispatch = useDispatch();
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();
  console.log('logout');
  useEffect(() => {
    console.log('[Logout] useEffect');
    setTimeout(() => {
      console.log('[Logout] useeffect timeout');
      dispatch(logout());
      dispatch(setDefaultMemberNetwork());
      removeCookie('access_token');
      localStorage.removeItem('trs_stepper');
      localStorage.removeItem('persist:root');
      // clear the store
      Store.dispatch({ type: 'RESET' });
      setIsLoading(false);
      navigate('/');
    }, 2000);
  }, [dispatch, navigate]);

  return (
    <Flex height={'100vh'} alignItems="center" justifyContent={'center'}>
      <Stack textAlign="center" py={20}>
        {isLoading && <Loader text="Logout ... " />}
      </Stack>
    </Flex>
  );
};

export default Logout;
