import { Box, Stack, Button, Text, useColorModeValue } from '@chakra-ui/react';

import { GoogleIcon } from 'components/Icon';

import { Trans } from '@lingui/react';
import ChakraRouterLink from 'components/ChakraRouterLink';
import LoginForm from './LoginForm';
import AuthLayout from 'layouts/AuthLayout';
interface LoginProps {
  handleSignWithSocial: (event: React.FormEvent, type: string) => void;
  handleSignWithEmail: (data: any) => void;
  isLoading?: boolean;
  isError?: any;
}

const Login: React.FC<LoginProps> = (props) => (
  <AuthLayout>
    <Text fontWeight="bold" color={useColorModeValue('gray.600', 'white')}>
      <Trans id="Log into your TRISA account.">Log into your TRISA account.</Trans>
    </Text>
    <Box>
      <Button
        data-testid="signin-with-google"
        bg={'gray.100'}
        w="100%"
        size="lg"
        borderRadius="none"
        onClick={(event: any) => props.handleSignWithSocial(event, 'google')}
        color={'gray.600'}
        _hover={{
          background: useColorModeValue('gray.200', 'black'),
          color: useColorModeValue('gray.600', 'white')
        }}
        _focus={{
          borderColor: 'transparent'
        }}>
        <GoogleIcon h={24} />
        <Text as={'span'} ml={3}>
          <Trans id="Continue with Google">Continue with Google</Trans>
        </Text>
      </Button>
    </Box>
    <Text align="center">or</Text>

    <Box bg={useColorModeValue('white', 'transparent')}>
      <LoginForm handleSignWithEmail={props.handleSignWithEmail} isLoading={props.isLoading} />
      <Text textAlign="center" fontSize="1rem">
        <Trans id="Not a TRISA Member?">Not a TRISA Member?</Trans>{' '}
        <ChakraRouterLink to="/auth/register" color={'#1F4CED'}>
          <Trans id="Join the TRISA network today.">Join the TRISA network today.</Trans>
        </ChakraRouterLink>
      </Text>
    </Box>
  </AuthLayout>
);

export default Login;
