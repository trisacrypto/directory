import { Box, Button, Heading, VStack, HStack, Text, useColorModeValue } from '@chakra-ui/react';

import { GoogleIcon, GithubIcon, MicrosoftIcon } from 'components/Icon';

import { Trans } from '@lingui/react';
import ChakraRouterLink from 'components/ChakraRouterLink';
import LoginForm from 'components/Form/LoginForm';
import AuthLayout from 'layouts/AuthLayout';
interface LoginProps {
  handleSignWithSocial: (event: React.FormEvent, type: string) => void;
  handleSignWithEmail: (data: any) => void;
  isLoading?: boolean;
  isError?: any;
}

const Login: React.FC<LoginProps> = ({ handleSignWithSocial, handleSignWithEmail, isLoading }) => (
  <AuthLayout>
    <Heading
      fontWeight="bold"
      color={useColorModeValue('gray.600', 'white')}
      textTransform="capitalize"
      textAlign="center"
      size="md">
      <Trans id="Log into your TRISA account">Log into your TRISA account</Trans>
    </Heading>
    <VStack>
      <Button
        data-testid="signin-with-google"
        bg={'white'}
        border="1px gray solid"
        pl={-2}
        w="100%"
        size="lg"
        borderRadius="none"
        onClick={(event: any) => handleSignWithSocial(event, 'google')}
        color={'gray.600'}
        _hover={{
          background: useColorModeValue('gray.200', 'black'),
          color: useColorModeValue('gray.600', 'white')
        }}
        _focus={{
          borderColor: 'transparent'
        }}>
        <HStack spacing={5}>
          <Box pos={'absolute'} left={5}>
            <GoogleIcon h={24} />
          </Box>
          <Text as={'span'} fontSize="md">
            <Trans id="Continue with Google">Continue with Google</Trans>
          </Text>
        </HStack>
      </Button>
      <Button
        data-testid="signin-with-github"
        bg={'white'}
        border="1px gray solid"
        w="100%"
        pl={-2}
        size="lg"
        borderRadius="none"
        onClick={(event: any) => handleSignWithSocial(event, 'github')}
        color={'gray.600'}
        _hover={{
          background: useColorModeValue('gray.200', 'black'),
          color: useColorModeValue('gray.600', 'white')
        }}
        _focus={{
          borderColor: 'transparent'
        }}>
        <HStack spacing={5}>
          <Box pos={'absolute'} left={5}>
            <GithubIcon h={24} />
          </Box>
          <Text as={'span'} fontSize="md">
            <Trans id="Continue with GitHub">Continue with GitHub</Trans>
          </Text>
        </HStack>
      </Button>
      <Button
        data-testid="signin-with-microsoft"
        bg={'white'}
        border="1px gray solid"
        w="100%"
        size="lg"
        pl="4"
        borderRadius="none"
        onClick={(event: any) => handleSignWithSocial(event, 'microsoft')}
        color={'gray.600'}
        _hover={{
          background: useColorModeValue('gray.200', 'black'),
          color: useColorModeValue('gray.600', 'white')
        }}
        _focus={{
          borderColor: 'transparent'
        }}>
        <HStack spacing={5} justifyContent={'space-between'}>
          <Box pos={'absolute'} left={5}>
            <MicrosoftIcon h={24} />
          </Box>
          <Box as={'span'} fontSize="md">
            <Trans id="Continue with Microsoft">Continue with Microsoft</Trans>
          </Box>
        </HStack>
      </Button>
    </VStack>
    <Text align="center">or</Text>

    <Box bg={useColorModeValue('white', 'transparent')}>
      <LoginForm handleSignWithEmail={handleSignWithEmail} isLoading={isLoading} />
      <Text textAlign="center" fontSize="md">
        <Trans id="Not a TRISA Member?">Not a TRISA Member?</Trans>{' '}
        <ChakraRouterLink
          to="/auth/register"
          color={'#1F4CED'}
          fontWeight={500}
          _hover={{ textDecor: 'underline' }}>
          <Trans id="Join the TRISA network today.">Join the TRISA network today.</Trans>
        </ChakraRouterLink>
      </Text>
    </Box>
  </AuthLayout>
);

export default Login;
