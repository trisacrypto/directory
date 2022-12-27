import { Trans } from '@lingui/macro';
import { HStack, IconButton, useToast, Button } from '@chakra-ui/react';
import { ProfileBlock } from './ProfileBlock';
import InputFormControl from '../ui/InputFormControl';
import { auth0ResetPassword } from 'utils/auth0.helper';
// import ChangePasswordModal from './ChangePasswordModal';
import React from 'react';
import { useSelector } from 'react-redux';
import { userSelector } from 'modules/auth/login/user.slice';
import CkLazyLoadImage from 'components/LazyImage';
import EditIcon from 'assets/edit-input.svg';
const PasswordInput = (props: any) => {
  return (
    <HStack w="100%" align="start">
      <InputFormControl {...props} />
      <IconButton
        aria-label="Edit"
        icon={<CkLazyLoadImage src={EditIcon} mx="auto" w="25px" />}
        variant="unstyled"
        marginTop="32px!important"
        onClick={(e) => props.updatePasswordHandler(e)}
      />
    </HStack>
  );
};
export const UserProfilePassword = () => {
  const toast = useToast();
  const { user } = useSelector(userSelector);
  const [isLoading, setIsLoading] = React.useState(false);

  const updatePasswordHandler = async (e: any) => {
    e.preventDefault();
    setIsLoading(true);
    const res = await auth0ResetPassword({
      email: user?.email,
      connection: user?.authType === 'auth0' ? 'Username-Password-Authentication' : 'google-oauth2'
    });
    if (res) {
      toast({
        title: 'Password reset email sent',
        description: 'Please check your email to reset your password',
        status: 'success',
        position: 'top-right',
        duration: 9000,
        isClosable: true
      });
    }
    setIsLoading(false);
  };

  return (
    <ProfileBlock title={<Trans>SECURITY</Trans>}>
      <Button
        color="white"
        fontWeight={700}
        fontSize="md"
        border={1}
        isLoading={isLoading}
        onClick={updatePasswordHandler}>
        Change password
      </Button>
    </ProfileBlock>
  );
};

export default React.memo(PasswordInput);
