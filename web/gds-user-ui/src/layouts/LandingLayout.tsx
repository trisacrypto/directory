import React, { useEffect } from 'react';
import { Box, Stack, Flex } from '@chakra-ui/react';
import LandingHeader from 'components/Header/LandingHeader';
import Footer from 'components/Footer/LandingFooter';
import useAuth from 'hooks/useAuth';
import { APP_PATH } from 'utils/constants';
import { useNavigate, useLocation } from 'react-router-dom';
type LandingLayoutProp = {
  children?: React.ReactNode;
};

export default function LandingLayout(props: LandingLayoutProp): JSX.Element {
  const location = useLocation();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const isLoggedIn = isAuthenticated();

  // if user is logged for login pathname redirect to dashboard
  useEffect(() => {
    if (
      isLoggedIn &&
      (location.pathname === APP_PATH.LOGIN || location.pathname === APP_PATH.REGISTER)
    ) {
      navigate('/dashboard/overview');
    }
  }, [isLoggedIn, location.pathname, navigate]);

  return (
    <Stack
      // align="center"
      justifyContent="space-between"
      minW={'100%'}
      m="0 auto"
      spacing={0}
      fontFamily={'Open Sans'}
      position={'relative'}
      minHeight={'100vh'}>
      <LandingHeader />
      <Stack flexGrow={1}>{props.children}</Stack>
      <Footer />
    </Stack>
  );
}
