import { Stack, Button, useColorModeValue, Box } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import InputFormControl from 'components/ui/InputFormControl';
import { SubmitHandler, useFormContext } from 'react-hook-form';
import { getValueByPathname } from 'utils/utils';
import { ResetPasswordFormValues } from './PasswordReset';

type PasswordResetFormProps = {
  onSubmit: SubmitHandler<ResetPasswordFormValues>;
};

function PasswordResetForm({ onSubmit }: PasswordResetFormProps) {
  const { register, formState, handleSubmit } = useFormContext<ResetPasswordFormValues>();
  const { isSubmitting, errors } = formState;

  return (
    <Box
      rounded={'lg'}
      bg={useColorModeValue('white', 'transparent')}
      position={'relative'}
      bottom={5}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Stack spacing={4}>
          <InputFormControl
            placeholder={t`Email Address`}
            type="email"
            size="lg"
            label="Enter your email address"
            isInvalid={getValueByPathname(errors, 'email')}
            formHelperText={getValueByPathname(errors, 'email')?.message}
            controlId={''}
            {...register('email')}
          />
          <Button
            display="block"
            alignSelf="start"
            px={16}
            bg="blue"
            color={'white'}
            isLoading={isSubmitting}
            type="submit"
            _hover={{
              bg: '#10aaed'
            }}
            _focus={{
              borderColor: 'transparent'
            }}>
            <Trans id="Submit">Submit</Trans>
          </Button>
        </Stack>
      </form>
    </Box>
  );
}

export default PasswordResetForm;
