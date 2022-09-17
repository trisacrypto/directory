import React from 'react';
import { Box, Stack, Button, Text, useColorModeValue } from '@chakra-ui/react';

import { GoogleIcon } from 'components/Icon';
import { useForm } from 'react-hook-form';

import { yupResolver } from '@hookform/resolvers/yup';
import { getValueByPathname } from 'utils/utils';
import InputFormControl from 'components/ui/InputFormControl';
import PasswordStrength from 'components/PasswordStrength';
import * as yup from 'yup';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import ChakraRouterLink from 'components/ChakraRouterLink';
import AuthLayout from 'layouts/AuthLayout';
import SignupForm from 'components/Form/SignupForm';
interface CreateAccountProps {
  handleSocialAuth: (event: React.FormEvent, type: string) => void;
  handleSignUpWithEmail: (data: any) => void;
  isLoading?: boolean;
  isError?: any;
  isPasswordError?: boolean;
  isUsernameError?: boolean;
}
interface IFormInputs {
  username: string;
  password: string;
}

const validationSchema = yup.object().shape({
  username: yup.string().email('Email is not valid').required('Email is required'),
  password: yup.string().required('Password is required')
});

// TO-DO : need some improvements
const CreateAccount: React.FC<CreateAccountProps> = (props) => {
  return (
    <AuthLayout>
      <Text color={useColorModeValue('gray.600', 'white')}>
        <Text as={'span'} fontWeight={'bold'}>
          <Trans id="Create your TRISA account.">Create your TRISA account.</Trans>
        </Text>{' '}
        <Trans id="We recommend that a senior compliance officer initially creates the account for the VASP. Additional accounts can be created later.">
          We recommend that a senior compliance officer initially creates the account for the VASP.
          Additional accounts can be created later.
        </Trans>
      </Text>

      <Box>
        <Button
          bg={'gray.100'}
          w="100%"
          onClick={(event: any) => props.handleSocialAuth(event, 'google')}
          size="lg"
          borderRadius="none"
          color={'gray.600'}
          _hover={{
            background: useColorModeValue('gray.200', 'black'),
            color: useColorModeValue('gray.600', 'white')
          }}>
          <GoogleIcon h={24} />
          <Text as={'span'} ml={3}>
            <Trans id="Continue with Google">Continue with Google</Trans>
          </Text>
        </Button>
      </Box>
      <Text textAlign="center">Or</Text>
      <Box bg={useColorModeValue('white', 'transparent')}>
        <SignupForm
          handleSignUpWithEmail={props.handleSignUpWithEmail}
          isLoading={props.isLoading}
        />
      </Box>
    </AuthLayout>
  );
};

export default CreateAccount;
