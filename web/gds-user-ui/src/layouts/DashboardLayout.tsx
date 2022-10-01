import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import Sidebar from 'components/Sidebar';
import { Box } from '@chakra-ui/react';
import Loader from 'components/Loader';
import { userSelector } from 'modules/auth/login/user.slice';
import TransparentLoader from '../components/Loader/TransparentLoader';
import { isTokenExpired } from 'utils/utils';
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
    // check if token is valid
    // if not, redirect to login page
    if (isTokenExpired()) {
      setIsLoading(true);
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
