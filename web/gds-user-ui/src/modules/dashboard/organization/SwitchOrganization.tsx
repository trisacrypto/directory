import React, { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
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
  const query = useQuery();
  const vaspName = query.get('vaspName');
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const isCalled = useRef(false);
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
          const token = (await refreshNewToken()) as any;
          const user = token && (await getUserCurrentOrganizationService());
          if (user?.status === APP_STATUS_CODE.OK) {
            dispatch(setUserOrganization(user?.data));
            setIsLoading(false);
            navigate(APP_PATH.DASHBOARD);
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
      } finally {
        setIsLoading(false);
      }
    };

    if (!isCalled.current) {
      switchOrganization();
      isCalled.current = true;
    }

    return () => {
      isCalled.current = false;
    };
  }, [id, navigate, dispatch, toast, vaspName]);

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

  return <>{isLoading && <TransparentLoader title={renderTitle()} opacity="full" />}</>;
};

export default SwitchOrganization;
