import { ReactNode } from 'react';
import { FormLabel, Heading, HStack, Stack, Text, VStack } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import InputFormControl, { _FormControlProps } from 'components/ui/InputFormControl';
import { SimpleDashboardLayout } from 'layouts';
import FormLayout from 'layouts/FormLayout';
import UserProfileIcon from 'assets/ph_user-circle-plus-light.svg';
import CkLazyLoadImage from 'components/LazyImage';
import AddLinkedAccountModal from './AddLinkedAccountModal';
import RemoveLinkedAccountModal from './RemoveLinkedAccountModal';
import ChangeNameModal from './ChangeNameModal';
import ChangePasswordModal from './ChangePasswordModal';

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

const PasswordInput = (props: _FormControlProps) => {
  return (
    <HStack w="100%" align="start">
      <InputFormControl {...props} />
      <ChangePasswordModal />
    </HStack>
  );
};

function UserProfile() {
  return (
    <SimpleDashboardLayout>
      <Heading size="lg" mb={5}>
        <Trans>User Profile</Trans>
      </Heading>
      <FormLayout>
        <VStack w="100%" align="start" spacing={8}>
          <ProfileBlock title={<Trans>Login & Identity</Trans>}>
            <Stack direction="row" justifyContent="space-between" w="100%">
              <VStack align="start">
                <VStack align="start">
                  <Text fontWeight={700}>
                    <Trans>Email Address</Trans>
                  </Text>
                  <Text>jferdinand@vaspnet.co.uk</Text>
                </VStack>
                <VStack align="start">
                  <Text fontWeight={700}>
                    <Trans>Account ID</Trans>
                  </Text>
                  <Text>0087765</Text>
                </VStack>
              </VStack>
              <VStack>
                <CkLazyLoadImage src={UserProfileIcon} mx="auto" />
              </VStack>
            </Stack>

            <EditableInput
              label={
                <FormLabel fontWeight={700}>
                  <Trans>First (Given) Name</Trans>
                </FormLabel>
              }
              controlId="first_given_name"
            />
            <EditableInput
              label={
                <FormLabel fontWeight={700}>
                  <Trans>Last (Family) Name</Trans>
                </FormLabel>
              }
              controlId="first_given_name"
            />
          </ProfileBlock>

          <ProfileBlock title={<Trans>SECURITY</Trans>}>
            <PasswordInput
              type="password"
              value="blablablabla"
              label={
                <FormLabel fontWeight={700}>
                  <Trans>Password</Trans>
                </FormLabel>
              }
              controlId="password"
            />
          </ProfileBlock>

          <ProfileBlock title={<Trans>USER DETAILS</Trans>}>
            <VStack align="start">
              <Text fontWeight={700} textTransform="capitalize">
                <Trans>Profile Created</Trans>
              </Text>
              <Text>January 3, 2022</Text>
            </VStack>
            <VStack align="start">
              <Text fontWeight={700} textTransform="capitalize">
                <Trans>Role</Trans>
              </Text>
              <Text textTransform="capitalize">
                <Trans>Organization Leader</Trans>
              </Text>
            </VStack>
            <VStack align="start">
              <Text fontWeight={700} textTransform="capitalize">
                <Trans>Permissions</Trans>
              </Text>
              <Text>
                <Trans>
                  Create new organization, add/approve collaborators, submit certificate request,
                  check status of certificate request
                </Trans>
              </Text>
            </VStack>
            <VStack align="start">
              <Text fontWeight={700} textTransform="capitalize">
                <Trans>Last Login</Trans>
              </Text>
              <Text>March 4, 2022</Text>
            </VStack>
          </ProfileBlock>

          <ProfileBlock
            title={
              <>
                <Trans>LINKED ACCOUNTS</Trans>
                <AddLinkedAccountModal />
              </>
            }>
            <Text>
              <Trans>
                If you have additional accounts with the TRISA Global Directory Service, you can
                link them here. You will be required to log in to the linked account to verify
                account ownership.
              </Trans>
            </Text>
            <HStack w="100%">
              <InputFormControl
                label={
                  <FormLabel fontWeight={700}>
                    <Trans>Linked Account</Trans>
                  </FormLabel>
                }
                controlId="linked_account"
                placeholder="sdze"
              />
              <RemoveLinkedAccountModal />
            </HStack>
          </ProfileBlock>
        </VStack>
      </FormLayout>
    </SimpleDashboardLayout>
  );
}

export default UserProfile;
