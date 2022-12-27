import React from 'react';
import { useParams, Navigate } from 'react-router-dom';
import useQuery from 'hooks/useQuery';
import TransparentLoader from 'components/Loader/TransparentLoader';

import { useToast, Text } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useSwitchOrganization } from './useSwitchOrganization';
import { userSelector } from 'modules/auth/login/user.slice';
import { useSelector } from 'react-redux';
import { APP_PATH } from 'utils/constants';
const SwitchOrganization: React.FC = () => {
  const toast = useToast();
  const { id } = useParams<{ id: string }>() as any;
  const { vaspName } = useQuery<{ vaspName: string }>();
  const { isLoading, isError } = useSwitchOrganization(id);
  const { user } = useSelector(userSelector);

  if (isError) {
    toast({
      title: 'Organization not found',
      status: 'error',
      duration: 5000,
      isClosable: true,
      position: 'top-right'
    });
    return <Navigate to={APP_PATH.DASHBOARD} />;
  }

  const renderLoadingTitle = () => {
    return (
      <Text as={'span'}>
        Switching to{' '}
        <Text as={'span'} color={colors.system.blue} fontWeight={'bold'}>
          {vaspName || user?.vaps?.name} ...
        </Text>
      </Text>
    );
  };

  return (
    <>
      {isLoading && !isError ? (
        <TransparentLoader title={renderLoadingTitle()} opacity="full" />
      ) : (
        <Navigate to={APP_PATH.DASHBOARD} />
      )}
    </>
  );
};

export default SwitchOrganization;
