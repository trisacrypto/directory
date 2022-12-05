import { Trans } from '@lingui/macro';
import { Text, VStack, Tag } from '@chakra-ui/react';
import { ProfileBlock } from './index';
import { useSelector } from 'react-redux';
import { userSelector } from 'modules/auth/login/user.slice';

const UserDetails: React.FC = () => {
  const { user } = useSelector(userSelector);

  return (
    <>
      <ProfileBlock title={<Trans>USER DETAILS</Trans>}>
        <VStack align="start">
          <Text fontWeight={700} textTransform="capitalize">
            <Trans>Profile Created</Trans>
          </Text>
          <Text data-testid="user_created_At">{user?.createAt || '-'}</Text>
        </VStack>
        <VStack align="start">
          <Text fontWeight={700} textTransform="capitalize">
            <Trans>Role</Trans>
          </Text>
          <Text textTransform="capitalize" data-testid="user_role">
            <Trans>{user?.role || '-'}</Trans>
          </Text>
        </VStack>
        <VStack align="start">
          <Text fontWeight={700} textTransform="capitalize">
            <Trans>Permissions</Trans>
          </Text>
          <Text>
            {user?.permissions.map((permission: string, index: string) => (
              <Tag key={index} bg={'blue'} color={'white'} ml={1} data-testid={`user_permissions`}>
                {permission}
              </Tag>
            ))}
          </Text>
        </VStack>
        <VStack align="start">
          <Text fontWeight={700} textTransform="capitalize">
            <Trans>Last Login</Trans>
          </Text>
          <Text data-testid="user_last_login">{user?.lastLogin}</Text>
        </VStack>
      </ProfileBlock>
    </>
  );
};

export default UserDetails;
