import { Button, FormLabel, VStack } from '@chakra-ui/react';
import { FC } from 'react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import InputFormControl from 'components/ui/InputFormControl';
import { Trans } from '@lingui/macro';
import * as yup from 'yup';
// import useCreateUserName  from 'hooks/useCreateUserName';
type Props = {
  onCloseModal: () => void;
};

type TName = {
  first_name: string;
  last_name: string;
};

const ChangeNameForm: FC<Props> = (props) => {
  const { onCloseModal } = props;

  const {
    handleSubmit,
    register,
    formState: { errors }
  } = useForm<TName>({
    mode: 'onBlur',
    resolver: yupResolver(
      yup.object().shape({
        first_name: yup.string().required('First name is required'),
        last_name: yup.string().required('Last name is required')
      })
    )
  });

  const onSubmit = (data: TName) => {
    console.log('submit', data);
  };

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <VStack align="start">
          <InputFormControl
            label={
              <FormLabel fontWeight={700}>
                <Trans>First (Given) Name</Trans>
              </FormLabel>
            }
            {...register('first_name')}
            data-testid="first_name"
            isInvalid={!!errors.first_name}
            formHelperText={errors.first_name?.message}
            controlId="first_name"
          />
          <InputFormControl
            label={
              <FormLabel fontWeight={700}>
                <Trans>Last (Family) Name</Trans>
              </FormLabel>
            }
            {...register('last_name')}
            isInvalid={!!errors.last_name}
            data-testid="last_name"
            formHelperText={errors.last_name?.message}
            controlId="last_name"
          />
        </VStack>

        <VStack display="flex" flexDir="column" rowGap={2}>
          <Button bg="orange" _hover={{ bg: 'orange' }} minW="150px" type="submit">
            Save
          </Button>
          <Button variant="ghost" onClick={onCloseModal}>
            Cancel
          </Button>
        </VStack>
      </form>
    </>
  );
};

export default ChangeNameForm;
