import React, { useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import Sidebar from 'components/Sidebar';
import { Box, Flex, Text, Spinner } from '@chakra-ui/react';
import { getCookie } from 'utils/cookies';
import Loader from 'components/Loader';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import { useNavigate } from 'react-router-dom';
import { getAuth0User, userSelector } from 'modules/auth/login/user.slice';
import { getRegistrationData } from 'modules/dashboard/registration/service';
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
  const { isFetching, isLoggedIn } = useSelector(userSelector);

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
