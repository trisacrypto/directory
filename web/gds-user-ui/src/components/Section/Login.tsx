import {
  Flex,
  Box,
  FormControl,
  Input,
  Stack,
  Link,
  Button,
  Heading,
  Text,
  useColorModeValue
} from '@chakra-ui/react';
import * as yup from 'yup';

import { GoogleIcon } from 'components/Icon';

import { colors } from 'utils/theme';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import InputFormControl from 'components/ui/InputFormControl';
import { getValueByPathname } from 'utils/utils';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import ChakraRouterLink from 'components/ChakraRouterLink';
interface LoginProps {
  handleSignWithSocial: (event: React.FormEvent, type: string) => void;
  handleSignWithEmail: (data: any) => void;
  isLoading?: boolean;
  isError?: any;
}
interface IFormInputs {
  username: string;
  password: string;
}

const defaultValues = {
  username: '',
  password: ''
};

const validationSchema = yup.object().shape({
  username: yup.string().email('Email Address is not valid').required('Email Address is required'),
  password: yup.string().required('Password is required')
});

const Login: React.FC<LoginProps> = (props) => {
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<IFormInputs>({ resolver: yupResolver(validationSchema), defaultValues });

  return (
    <Flex fontSize={'xl'} bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={8} mx={'auto'} maxW={'xl'} py={12} width={'100%'}>
        <Text fontSize={'xl'} fontWeight="bold" color={useColorModeValue('gray.600', 'white')}>
          <Trans id="Log into your TRISA account.">Log into your TRISA account.</Trans>
        </Text>
        <Stack align={'center'} justify={'center'} fontFamily={colors.font}>
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
              Continue with Google
            </Text>
          </Button>
        </Stack>
        <Text py={3} align="center">
          or
        </Text>

        <Box bg={useColorModeValue('white', 'transparent')}>
          <form onSubmit={handleSubmit(props.handleSignWithEmail)} noValidate>
            <Stack spacing={4}>
              <InputFormControl
                data-testid="email"
                controlId="email"
                placeholder={t`Email Address`}
                type="email"
                size="lg"
                isInvalid={getValueByPathname(errors, 'username')}
                formHelperText={getValueByPathname(errors, 'username')?.message}
                {...register('username')}
              />
              <InputFormControl
                data-testid="password"
                controlId="password"
                size="lg"
                placeholder={t`Password`}
                type="password"
                isInvalid={getValueByPathname(errors, 'password')}
                formHelperText={getValueByPathname(errors, 'password')?.message}
                {...register('password')}
              />
              <Stack direction={['column', 'row']} py="5" justifyContent="space-between">
                <Button
                  data-testid="login-btn"
                  bg={colors.system.blue}
                  color={'white'}
                  isLoading={props.isLoading}
                  px={2}
                  py={4}
                  w={['full', '50%']}
                  type="submit"
                  _hover={{
                    bg: '#10aaed'
                  }}
                  _focus={{
                    borderColor: 'transparent'
                  }}>
                  <Trans id="Log In">Log In</Trans>
                </Button>
                <Text display="flex" alignItems="flex-end" style={{ marginRight: '2rem' }}>
                  <ChakraRouterLink
                    to={'/auth/reset'}
                    color="#1F4CED"
                    fontFamily="Open sans, sans-serif"
                    fontSize="1rem">
                    <Trans id="Forgot password?">Forgot password?</Trans>
                  </ChakraRouterLink>
                </Text>
              </Stack>
            </Stack>
          </form>
          <Text textAlign="center" fontSize="1rem">
            <Trans id="Not a TRISA Member?">Not a TRISA Member?</Trans>{' '}
            <ChakraRouterLink to="/auth/register" color={'#1F4CED'}>
              <Trans id="Join the TRISA network today.">Join the TRISA network today.</Trans>
            </ChakraRouterLink>
          </Text>
        </Box>
      </Stack>
    </Flex>
  );
};
export default Login;
