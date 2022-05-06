import React, { useEffect, useState } from 'react';
import LandingLayout from 'layouts/LandingLayout';
import { Heading, Stack, Spinner, Flex, Box } from '@chakra-ui/react';
import useHashQuery from 'hooks/useHashQuery';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import Cookies from 'universal-cookie';
import AlertMessage from 'components/ui/AlertMessage';
import { useNavigate } from 'react-router-dom';

const Logout: React.FC = () => {
  const cookies = new Cookies();
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();
  useEffect(() => {
    setTimeout(() => {
      cookies.remove('access_token', { path: '/' });
      setIsLoading(false);
      navigate('/');
    }, 2000);
  });

  return (
    <Flex height={'100vh'} alignItems="center" justifyContent={'center'}>
      <Stack textAlign="center" py={20}>
        {isLoading && <Spinner size="xl" speed="0.65s" />}
      </Stack>
    </Flex>
  );
};

export default Logout;
