import React from 'react';
import { useSelector } from 'react-redux';
import Sidebar from 'components/Sidebar';
import { Box } from '@chakra-ui/react';
import Loader from 'components/Loader';
import { userSelector } from 'modules/auth/login/user.slice';
import TransparentLoader from '../components/Loader/TransparentLoader';
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

  return (
    <>
      {isFetching ? (
        <Loader />
      ) : (
        <>
          <AxiosErrorLoader />
          <Sidebar {...props} />
        </>
      )}
    </>
  );
};
export default DashboardLayout;
