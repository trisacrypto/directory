import React from 'react';
import { Box, Button, Heading, Text, useColorModeValue } from '@chakra-ui/react';

import { GoogleIcon } from 'components/Icon';
import { Trans } from '@lingui/react';
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

// TO-DO : need some improvements
const CreateAccount: React.FC<CreateAccountProps> = (props) => {
  return (
    <AuthLayout>
      <Heading
        fontWeight={'bold'}
        color={useColorModeValue('gray.600', 'white')}
        size="md"
        textAlign="center"
        textTransform="capitalize">
        <Trans id="Create your TRISA account">Create your TRISA account</Trans>
      </Heading>
      <Text color={useColorModeValue('gray.600', 'white')} mt="4px!important" fontSize="md">
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
          <Text as={'span'} ml={3} fontSize="md">
            <Trans id="Continue with Google">Continue with Google</Trans>
          </Text>
        </Button>
      </Box>
      <Text textAlign="center">or</Text>
      <Box
        color={useColorModeValue('gray.600', 'white')}
        bg={useColorModeValue('white', 'transparent')}>
        <SignupForm
          handleSignUpWithEmail={props.handleSignUpWithEmail}
          isLoading={props.isLoading}
        />
      </Box>
    </AuthLayout>
  );
};

export default CreateAccount;
