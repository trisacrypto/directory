import { Button, FormLabel, VStack, useToast } from '@chakra-ui/react';
import { FC, useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import InputFormControl from 'components/ui/InputFormControl';
import { Trans } from '@lingui/macro';
import * as yup from 'yup';
import { useCreateFullName } from './useCreateUserName';
import { setUserName, userSelector } from 'modules/auth/login/user.slice';
import { useSelector, useDispatch } from 'react-redux';

type Props = {
  onCloseModal: () => void;
};

type TName = {
  first_name: string;
  last_name: string;
};

const ChangeNameForm: FC<Props> = (props) => {
  const toast = useToast();
  const dispatch = useDispatch();
  const { user } = useSelector(userSelector);
  const { onCloseModal } = props;
  const { updateName, isUpdating, wasUpdated, hasUpdateFailed } = useCreateFullName();
  const [userFullName, setUserFullName] = useState<string>('');

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
    const fullname = `${data.first_name} ${data.last_name}`;
    updateName(fullname);

    setUserFullName(fullname);
  };

  useEffect(() => {
    if (wasUpdated) {
      dispatch(setUserName(userFullName));
      onCloseModal();
      toast({
        position: 'top-right',
        title: 'Name updated successfully',
        isClosable: true,
        status: 'success',
        duration: 9000
      });
    }
  }, [wasUpdated, onCloseModal, userFullName, toast, dispatch]);

  useEffect(() => {
    if (hasUpdateFailed) {
      onCloseModal();
      toast({
        position: 'top-right',
        title: 'Failed to update name',
        isClosable: true,
        status: 'error',
        duration: 9000
      });
    }
  }, [hasUpdateFailed, onCloseModal, toast]);

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <VStack align="start">
          <InputFormControl
            label={
              <FormLabel fontWeight={700}>
                <Trans>Current Name</Trans>
              </FormLabel>
            }
            value={user?.name}
            isDisabled
            data-testid="current_name"
            controlId="current_name"
          />
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
          <Button
            bg="orange"
            _hover={{ bg: 'orange' }}
            minW="150px"
            type="submit"
            data-testid="save_button"
            isLoading={isUpdating}
            disabled={isUpdating}>
            Save
          </Button>
          <Button variant="ghost" onClick={onCloseModal} isDisabled={isUpdating}>
            Cancel
          </Button>
        </VStack>
      </form>
    </>
  );
};

export default ChangeNameForm;
