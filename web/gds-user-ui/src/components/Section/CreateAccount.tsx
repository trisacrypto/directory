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
import { getValueByPathname } from 'utils/utils';
import InputFormControl from 'components/ui/InputFormControl';

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

// const validationSchema = yup.object().shape({
//   username: yup.string().email('Email is not valid').required('Email is required'),
//   password: yup.string().required('Password is required')
// });

// TO-DO : need some improvements
const CreateAccount: React.FC<CreateAccountProps> = (props) => {
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<IFormInputs>({
    resolver: yupResolver(validationSchema)
  });
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
          <form onSubmit={handleSubmit(props.handleSignUpWithEmail)} noValidate>
            <Stack spacing={4}>
              <InputFormControl
                controlId=""
                {...register('username')}
                paddingY={6}
                data-testid="username-field"
                placeholder="Email Address"
                isInvalid={!!getValueByPathname(errors, 'username')}
                formHelperText={getValueByPathname(errors, 'username')?.message}
              />

              <InputFormControl
                controlId=""
                {...register('password')}
                paddingY={6}
                data-testid="password-field"
                placeholder="Password"
                hasBtn
                handleFn={handleClick}
                setBtnName={show ? 'Hide' : 'Show'}
                isInvalid={!!getValueByPathname(errors, 'password')}
                type={show ? 'text' : 'password'}
                formHelperText={
                  getValueByPathname(errors, 'password') ? (
                    getValueByPathname(errors, 'password')?.message
                  ) : (
                    <>
                      * At least 8 characters in length * Contain at least 3 of the following 4
                      types of characters: * lower case letters (a-z) * upper case letters (A-Z) *
                      numbers (i.e. 0-9) * special characters (e.g. !@#$%^&*)
                    </>
                  )
                }
              />
              <Stack spacing={10}>
                <Button
                  bg={colors.system.blue}
                  color={'white'}
                  py={6}
                  type="submit"
                  borderRadius={'none'}
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
