import InputFormControl from 'components/ui/InputFormControl';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';
import { getValueByPathname } from 'utils/utils';
import { Button, chakra, Stack, Text } from '@chakra-ui/react';
import ChakraRouterLink from 'components/ChakraRouterLink';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';

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

type LoginFormProps = {
  handleSignWithEmail: (args: any) => void;
  isLoading?: boolean;
};

function LoginForm(props: LoginFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<IFormInputs>({ resolver: yupResolver(validationSchema), defaultValues });

  return (
    <chakra.form onSubmit={handleSubmit(props.handleSignWithEmail)} noValidate>
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
            color={'white'}
            isLoading={props.isLoading}
            px={2}
            py={4}
            w={['full', '50%']}
            type="submit">
            <Trans id="Log In">Log In</Trans>
          </Button>
          <Text display="flex" alignItems="flex-end" style={{ marginRight: '2rem' }}>
            <ChakraRouterLink to="/auth/reset" color="#1F4CED" fontSize="1rem">
              <Trans id="Forgot password?">Forgot password?</Trans>
            </ChakraRouterLink>
          </Text>
        </Stack>
      </Stack>
    </chakra.form>
  );
}

export default LoginForm;
