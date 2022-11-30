import React, { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import useQuery from 'hooks/useQuery';
import TransparentLoader from 'components/Loader/TransparentLoader';
import { logUserInBff } from 'modules/auth/login/auth.service';
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

  useEffect(() => {
    const switchOrganization = async () => {
      try {
        const logged = await logUserInBff({
          orgid: id
        });
        if (logged.status === APP_STATUS_CODE.NO_CONTENT) {
          await refreshNewToken();
          dispatch(
            setUserOrganization({
              organization: vaspName
            })
          );
          navigate(APP_PATH.DASHBOARD);
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

        console.log('error', error);
      }
    };

    if (!isCalled.current) {
      switchOrganization();
      isCalled.current = true;
    }

    return () => {
      isCalled.current = false;
    };
  }, [id, navigate, dispatch, toast, query, vaspName]);

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

  return <>{!isError && <TransparentLoader title={renderTitle()} opacity="full" />}</>;
};

export default SwitchOrganization;
