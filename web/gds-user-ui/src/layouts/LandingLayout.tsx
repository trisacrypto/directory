import React, { useEffect } from 'react';
import { Stack } from '@chakra-ui/react';
import LandingHeader from 'components/Header/LandingHeader';
import Footer from 'components/Footer/LandingFooter';
import useAuth from 'hooks/useAuth';
import { APP_PATH } from 'utils/constants';
import { useNavigate, useLocation } from 'react-router-dom';
import useSearchParams from 'hooks/useQueryParams';

type LandingLayoutProp = {
  children?: React.ReactNode;
};

export default function LandingLayout(props: LandingLayoutProp): JSX.Element {
  const location = useLocation();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const isLoggedIn = isAuthenticated();
  const { orgid } = useSearchParams();
  // if user is logged for login pathname redirect to dashboard
  useEffect(() => {
    if (
      isLoggedIn &&
      (location.pathname === APP_PATH.LOGIN || location.pathname === APP_PATH.REGISTER)
    ) {
      if (orgid) {
        const link = `${APP_PATH.SWITCH_ORGANIZATION}/${orgid}`;
        navigate(link);
        return;
      }
      navigate('/dashboard/overview');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isLoggedIn, location.pathname, orgid]);

  return (
    <Stack
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
