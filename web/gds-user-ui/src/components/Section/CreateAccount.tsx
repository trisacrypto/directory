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
// import { validationSchema } from 'modules/auth/register/register.validation';
import { getValueByPathname } from 'utils/utils';
import InputFormControl from 'components/ui/InputFormControl';
import PasswordStrength from 'components/PasswordStrength';
import * as yup from 'yup';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import { Link as RouterLink } from 'react-router-dom';

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
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch
  } = useForm<IFormInputs>({
    resolver: yupResolver(validationSchema)
  });
  const [show, setShow] = React.useState(false);
  const handleClick = () => setShow(!show);
  const watchPassword = watch('password');

  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={'open sans'}
      fontSize={'xl'}
      mb={'10vh'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={8} mx={'auto'} maxW={'xl'} py={12} px={6}>
        <Stack align={'center'}>
          <Heading fontSize={'4xl'}></Heading>
          <Text color={useColorModeValue('gray.600', 'white')}>
            <Text as={'span'} fontWeight={'bold'}>
              <Trans id="Create your TRISA account.">Create your TRISA account.</Trans>
            </Text>{' '}
            <Trans id="We recommend that a senior compliance officer initialally creates the account for the VASP. Additional accounts can be created later.">
              We recommend that a senior compliance officer initialally creates the account for the
              VASP. Additional accounts can be created later.
            </Trans>
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
              <Trans id="Continue with Google">Continue with Google</Trans>
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
                placeholder={t`Email Address`}
                isInvalid={!!getValueByPathname(errors, 'username')}
                formHelperText={getValueByPathname(errors, 'username')?.message}
              />

              <InputFormControl
                controlId=""
                {...register('password')}
                paddingY={6}
                data-testid="password-field"
                placeholder={t`Password`}
                hasBtn
                handleFn={handleClick}
                setBtnName={show ? 'Hide' : 'Show'}
                isInvalid={!!getValueByPathname(errors, 'password')}
                type={show ? 'text' : 'password'}
                formHelperText={watchPassword ? <PasswordStrength data={watchPassword} /> : null}
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
                  <Trans id="Create an Account">Create an Account</Trans>
                </Button>
                <Text textAlign="center">
                  <Trans id="Already have an account?">Already have an account?</Trans>{' '}
                  <RouterLink to={'/auth/login'}>
                    <Link color={colors.system.cyan}>
                      {' '}
                      <Trans id="Log in.">Log in.</Trans>
                    </Link>
                  </RouterLink>
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
