import React from 'react';
import { Box, Button, Heading, Text, useColorModeValue, VStack, HStack } from '@chakra-ui/react';
import { GoogleIcon, GithubIcon, MicrosoftIcon } from 'components/Icon';
import { Trans } from '@lingui/react';
import AuthLayout from 'layouts/AuthLayout';
import SignupForm from 'components/Form/SignupForm';
import LandingBanner from 'components/Banner/LandingBanner';
interface CreateAccountProps {
  handleSocialAuth: (event: React.FormEvent, type: string) => void;
  handleSignUpWithEmail: (data: any) => void;
  isLoading?: boolean;
  isError?: any;
  isPasswordError?: boolean;
  isUsernameError?: boolean;
}

// TO-DO : need some improvements
const CreateAccount: React.FC<CreateAccountProps> = ({
  handleSocialAuth,
  handleSignUpWithEmail,
  isLoading
}) => {
  return (
    <AuthLayout>
      <LandingBanner />
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

      <VStack>
        <Button
          data-testid="signin-with-google"
          bg={'white'}
          border="1px gray solid"
          pl={-2}
          w="100%"
          size="lg"
          borderRadius="none"
          onClick={(event: any) => handleSocialAuth(event, 'google')}
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
          onClick={(event: any) => handleSocialAuth(event, 'github')}
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
          onClick={(event: any) => handleSocialAuth(event, 'microsoft')}
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
      <Text textAlign="center">or</Text>
      <Box
        color={useColorModeValue('gray.600', 'white')}
        bg={useColorModeValue('white', 'transparent')}>
        <SignupForm handleSignUpWithEmail={handleSignUpWithEmail} isLoading={isLoading} />
      </Box>
    </AuthLayout>
  );
};

export default CreateAccount;
