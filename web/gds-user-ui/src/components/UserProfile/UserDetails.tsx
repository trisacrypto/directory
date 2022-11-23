import { Trans } from '@lingui/macro';
import { Text, VStack, Tag } from '@chakra-ui/react';
import { ProfileBlock } from './index';
import { useSelector } from 'react-redux';
import { userSelector } from 'modules/auth/login/user.slice';
export default function UserDetails() {
  const { user } = useSelector(userSelector);
  return (
    <ProfileBlock title={<Trans>USER DETAILS</Trans>}>
      <VStack align="start">
        <Text fontWeight={700} textTransform="capitalize">
          <Trans>Profile Created</Trans>
        </Text>
        <Text>{user?.createAt || '-'}</Text>
      </VStack>
      <VStack align="start">
        <Text fontWeight={700} textTransform="capitalize">
          <Trans>Role</Trans>
        </Text>
        <Text textTransform="capitalize">
          <Trans>{user?.role || '-'}</Trans>
        </Text>
      </VStack>
      <VStack align="start">
        <Text fontWeight={700} textTransform="capitalize">
          <Trans>Permissions</Trans>
        </Text>
        <Text>
          {user?.permissions.map((permission: string, index: string) => (
            <Tag key={index} bg={'blue'} color={'white'} ml={1}>
              {permission}
            </Tag>
          ))}
        </Text>
      </VStack>
      <VStack align="start">
        <Text fontWeight={700} textTransform="capitalize">
          <Trans>Last Login</Trans>
        </Text>
        <Text>{user?.lastLogin || '-'}</Text>
      </VStack>
    </ProfileBlock>
  );
}
