import React, { useEffect, useState } from 'react';
import { useParams, Navigate } from 'react-router-dom';
import useQuery from 'hooks/useQuery';
import TransparentLoader from 'components/Loader/TransparentLoader';
import { logUserInBff, getUserCurrentOrganizationService } from 'modules/auth/login/auth.service';
import { refreshNewToken } from 'utils/auth0.helper';
import { useDispatch } from 'react-redux';
import { useToast, Text } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { setUserOrganization } from 'modules/auth/login/user.slice';
// import { getAuth0User } from 'modules/auth/login/user.slice';
import { APP_PATH, APP_STATUS_CODE } from 'utils/constants';
const SwitchOrganization: React.FC = () => {
  const toast = useToast();
  const { id } = useParams<{ id: string }>();
  const { vaspName } = useQuery<{ vaspName: string }>();
  const dispatch = useDispatch();
  // const navigate = useNavigate();
  // const isCalled = useRef(false);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const switchOrganization = async () => {
      try {
        setIsLoading(true);
        const logged = await logUserInBff({
          orgid: id
        });
        if (logged.status === APP_STATUS_CODE.NO_CONTENT) {
          console.log('TEST switchOrganization -----', logged);
          const token = (await refreshNewToken()) as any;
          console.log('TEST token generate --------', token);
          const user = token && (await getUserCurrentOrganizationService());

          if (user?.status === APP_STATUS_CODE.OK) {
            console.log('[TEST user is there ----]', user);
            dispatch(setUserOrganization(user?.data));
            setIsLoading(false);
          }
        }
      } catch (error) {
        setIsError(true);
        toast({
          title: 'Error',
          description: 'Something went wrong or Organization not found',
          status: 'error',
          duration: 9000,
          isClosable: true
        });
      }
    };

    switchOrganization();

    // if (!isCalled.current) {
    //   switchOrganization();
    //   isCalled.current = true;
    // }

    // return () => {
    //   isCalled.current = false;
    // };
  }, [id]);

  const renderTitle = () => {
    if (isError) {
      return 'Error';
    }
    return (
      <Text as={'span'}>
        Switching to{' '}
        <Text as={'span'} color={colors.system.blue} fontWeight={'bold'}>
          {vaspName} ...
        </Text>
      </Text>
    );
  };

  return (
    <>
      {isLoading && !isError ? (
        <TransparentLoader title={renderTitle()} opacity="full" />
      ) : (
        <Navigate to={APP_PATH.DASHBOARD} />
      )}
    </>
  );
};

export default SwitchOrganization;
