import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import Sidebar from 'components/Sidebar';
import { Box } from '@chakra-ui/react';
import Loader from 'components/Loader';
import { userSelector } from 'modules/auth/login/user.slice';
import TransparentLoader from '../components/Loader/TransparentLoader';
import { isTokenExpired } from 'utils/utils';
import { clearCookies } from 'utils/cookies';
import { persistor } from 'application/store';
type DashboardLayoutProp = {
  children: React.ReactNode;
};
// add loading component that will be shown when axios error is thrown
const AxiosErrorLoader = () => {
  return (
    <>
      <Box id="axiosLoader" display={'none'}>
        <TransparentLoader />
      </Box>
    </>
  );
};

const DashboardLayout: React.FC<DashboardLayoutProp> = (props) => {
  const { isFetching } = useSelector(userSelector);
  const [isLoading, setIsLoading] = useState(false);

  // const deps = window.location.pathname;
  useEffect(() => {
    if (isTokenExpired()) {
      setIsLoading(true);
      clearCookies();
      persistor.purge();
      localStorage.removeItem('persist:root');
      localStorage.removeItem('trs_stepper');
      setTimeout(() => {
        window.location.href = `/auth/login?q=token_expired`;
      }, 1000);
    }
  }, []);

  return (
    <>
      {isFetching ? (
        <Loader />
      ) : (
        <>
          {isLoading && (
            <Box>
              <TransparentLoader title="Your session has expired..." />
            </Box>
          )}
          <AxiosErrorLoader />
          <Sidebar {...props} />
        </>
      )}
    </>
  );
};
export default DashboardLayout;
