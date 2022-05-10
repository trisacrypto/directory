import React, { useState } from 'react';
import {
  Flex,
  Box,
  FormControl,
  FormLabel,
  Input,
  Checkbox,
  Stack,
  Link,
  Button,
  Heading,
  Text,
  useColorModeValue,
  FormHelperText,
  FormErrorMessage,
  InputGroup,
  InputRightElement
} from '@chakra-ui/react';

import { GoogleIcon } from 'components/Icon';
import { colors } from 'utils/theme';
import { useForm } from 'react-hook-form';

import { yupResolver } from '@hookform/resolvers/yup';
import { validationSchema } from 'modules/auth/register/register.validation';
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
// TO-DO : need some improvements
const CreateAccount: React.FC<CreateAccountProps> = (props) => {
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<IFormInputs>();
  const [show, setShow] = React.useState(false);
  const handleClick = () => setShow(!show);

  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={'open sans'}
      fontSize={'xl'}
      marginTop={'10vh'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={8} mx={'auto'} maxW={'xl'} py={12} px={6}>
        <Stack align={'center'}>
          <Heading fontSize={'4xl'}></Heading>
          <Text color={useColorModeValue('gray.600', 'white')}>
            <Text as={'span'} fontWeight={'bold'}>
              Create your TRISA account.
            </Text>{' '}
            We recommend that a senior compliance officer initialally creates the account for the
            VASP. Additional accounts can be created later.
          </Text>
        </Stack>
        <Stack align={'center'} justify={'center'} fontFamily={'open sans'}>
          <Button
            bg={'gray.100'}
            w="100%"
            onClick={(event) => props.handleSocialAuth(event, 'google')}
            height={'64px'}
            color={'gray.600'}
            _hover={{
              background: useColorModeValue('gray.200', 'black'),
              color: useColorModeValue('gray.600', 'white')
            }}>
            <GoogleIcon h={24} />
            <Text as={'span'} ml={3}>
              Continue with google
            </Text>
          </Button>
          <Text py={3}>Or</Text>
        </Stack>

        <Box
          rounded={'lg'}
          bg={useColorModeValue('white', 'transparent')}
          width={'100%'}
          position={'relative'}
          bottom={5}>
          <form onSubmit={handleSubmit(props.handleSignUpWithEmail)}>
            <Stack spacing={4}>
              <FormControl id="email" isInvalid={props.isUsernameError}>
                <Input
                  type="email"
                  {...register('username')}
                  height={'64px'}
                  placeholder="Email Address"
                />
                {props.isUsernameError && (
                  <FormErrorMessage>
                    This username is already used , please change and retry,
                  </FormErrorMessage>
                )}
              </FormControl>
              <FormControl id="password" isInvalid={props.isPasswordError}>
                <InputGroup>
                  <Input
                    type={show ? 'text' : 'password'}
                    {...register('password')}
                    height={'64px'}
                    placeholder="Password"
                  />
                  <InputRightElement width="5.5rem" my={3}>
                    <Button h="2.75rem" size="md" onClick={handleClick}>
                      {show ? 'Hide' : 'Show'}
                    </Button>
                  </InputRightElement>
                </InputGroup>
                {!props.isPasswordError ? (
                  <FormHelperText>
                    * At least 8 characters in length * Contain at least 3 of the following 4 types
                    of characters: * lower case letters (a-z) * upper case letters (A-Z) * numbers
                    (i.e. 0-9) * special characters (e.g. !@#$%^&*)
                  </FormHelperText>
                ) : (
                  <FormErrorMessage>
                    * At least 8 characters in length * Contain at least 3 of the following 4 types
                    of characters: * lower case letters (a-z) * upper case letters (A-Z) *
                    numbers(i.e. 0-9) * special characters (e.g. !@#$%^&*)
                  </FormErrorMessage>
                )}
              </FormControl>
              <Stack spacing={10}>
                <Button
                  bg={colors.system.blue}
                  color={'white'}
                  height={'64px'}
                  type="submit"
                  isLoading={props.isLoading}
                  _hover={{
                    bg: '#10aaed'
                  }}>
                  Create an Account
                </Button>
                <Text textAlign="center">
                  Already have an account?{' '}
                  <Link href="/login" color={colors.system.cyan}>
                    {' '}
                    Log in.
                  </Link>
                </Text>
              </Stack>
            </Stack>
          </form>
        </Box>
      </Stack>
    </Flex>
  );
};

export default CreateAccount;
