import { ReactNode } from 'react';
import { FormLabel, Heading, HStack, Stack, Text, VStack } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import InputFormControl, { _FormControlProps } from 'components/ui/InputFormControl';
import FormLayout from 'layouts/FormLayout';
import UserProfileIcon from 'assets/ph_user-circle-plus-light.svg';
import CkLazyLoadImage from 'components/LazyImage';
import ChangeNameModal from './ChangeNameModal';
import { useSelector } from 'react-redux';
import { userSelector } from 'modules/auth/login/user.slice';
import UserDetails from './UserDetails';
import { UserProfilePassword } from './UserProfilePassword';

export const ProfileBlock = ({ title, children }: { title: ReactNode; children: ReactNode }) => {
  return (
    <VStack align="start" w="100%" spacing={5}>
      <Heading
        size="sm"
        textTransform="uppercase"
        display="flex"
        fontWeight={700}
        columnGap={4}
        alignItems="center"
        data-testid="profile_block_title">
        {title}
      </Heading>
      <VStack align="start" w="100%" spacing={4}>
        {children}
      </VStack>
    </VStack>
  );
};

const EditableInput = (props: _FormControlProps) => {
  return (
    <HStack w="100%" align="start">
      <InputFormControl {...props} />
      <ChangeNameModal />
    </HStack>
  );
};

function UserProfile() {
  const { user } = useSelector(userSelector);
  const isSocialConnection = () => user?.authType !== 'auth0';
  return (
    <>
      <Heading size="lg" mb={5}>
        <Trans>User Profile</Trans>
      </Heading>
      <FormLayout>
        <VStack w="100%" align="start" spacing={8}>
          <ProfileBlock title={<Trans>Login & Identity</Trans>}>
            <Stack direction="row" justifyContent="space-between" w="100%">
              <VStack align="start" spacing={3}>
                <div>
                  <Text fontWeight={700}>
                    <Trans>Email Address</Trans>
                  </Text>
                  <Text mt={'0 !important'}>{user?.email}</Text>
                </div>
                <div>
                  <Text fontWeight={700}>
                    <Trans>Account ID</Trans>
                  </Text>
                  <Text mt={'0 !important'}>{user?.id}</Text>
                </div>
                <div>
                  <Text fontWeight={700}>
                    <Trans>Provider</Trans>
                  </Text>
                  <Text mt={'0 !important'}>{user?.authType}</Text>
                </div>
              </VStack>
              <Stack>
                <CkLazyLoadImage
                  borderRadius="50%"
                  src={user?.pictureUrl || UserProfileIcon}
                  mx="auto"
                  h="150px"
                />
              </Stack>
            </Stack>

            <EditableInput
              label={
                <FormLabel fontWeight={700}>
                  <Trans>Full name</Trans>
                </FormLabel>
              }
              isDisabled={true}
              controlId="fullname"
              value={user?.name}
            />
          </ProfileBlock>

          {!isSocialConnection() && <UserProfilePassword />}

          <UserDetails />
        </VStack>
      </FormLayout>
    </>
  );
}

export default UserProfile;
