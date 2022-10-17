import React from 'react';
import InputFormControl from 'components/ui/InputFormControl';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';
import { getValueByPathname } from 'utils/utils';
import { Button, chakra, Stack, Text } from '@chakra-ui/react';
import ChakraRouterLink from 'components/ChakraRouterLink';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import PasswordStrength from 'components/PasswordStrength';
interface IFormInputs {
  username: string;
  password: string;
}

const validationSchema = yup.object().shape({
  username: yup.string().email('Email Address is not valid').required('Email Address is required'),
  password: yup.string().required('Password is required')
});

const defaultValues = {
  username: '',
  password: ''
};

type SignupFormProps = {
  handleSignUpWithEmail: (args: any) => void;
  isLoading?: boolean;
};

function SignupForm(props: SignupFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch
  } = useForm<IFormInputs>({ resolver: yupResolver(validationSchema), defaultValues });
  const [show, setShow] = React.useState(false);
  const handleClick = () => setShow(!show);
  const watchPassword = watch('password');
  return (
    <chakra.form onSubmit={handleSubmit(props.handleSignUpWithEmail)} noValidate>
      <Stack spacing={4}>
        <InputFormControl
          controlId="username"
          {...register('username')}
          size="lg"
          data-testid="username-field"
          placeholder={t`Email Address`}
          isInvalid={!!getValueByPathname(errors, 'username')}
          formHelperText={getValueByPathname(errors, 'username')?.message}
        />

        <InputFormControl
          controlId="password"
          {...register('password')}
          data-testid="password-field"
          placeholder={t`Password`}
          hasBtn
          size="lg"
          handleFn={handleClick}
          setBtnName={show ? 'Hide' : 'Show'}
          isInvalid={!!getValueByPathname(errors, 'password')}
          type={show ? 'text' : 'password'}
          formHelperText={watchPassword ? <PasswordStrength data={watchPassword} /> : null}
        />
        <Button
          display="block"
          alignSelf="center"
          bg={'blue'}
          color={'white'}
          type="submit"
          isLoading={props.isLoading}
          _hover={{
            bg: '#10aaed'
          }}
          fontSize="md">
          <Trans id="Create Your Account">Create Your Account</Trans>
        </Button>
        <Text textAlign="center" fontSize="md">
          <Trans id="Already have an account?">Already have an account?</Trans>{' '}
          <ChakraRouterLink
            to={'/auth/login'}
            color="link"
            fontWeight={500}
            _hover={{ textDecor: 'underline' }}>
            <Trans id="Log in.">Log in.</Trans>
          </ChakraRouterLink>
        </Text>
      </Stack>
    </chakra.form>
  );
}

export default SignupForm;
