import { FC, useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { ADD_COLLABORATOR_FORM_METHOD } from './AddCollaboratorFormValidation';
import type { NewCollaborator } from './AddCollaboratorType';
import { useCreateCollaborator } from 'hooks/useCreateCollaborator';
import { Button, Stack, Text, VStack, chakra, useToast } from '@chakra-ui/react';
import { DevTool } from '@hookform/devtools';
import { isProdEnv } from 'application/config';
import InputFormControl from 'components/ui/InputFormControl';
import CheckboxFormControl from 'components/ui/CheckboxFormControl';
import { t, Trans } from '@lingui/macro';
import { useFetchCollaborators } from 'components/Collaborators/useFetchCollaborator';
import { useSelector } from 'react-redux';
import { userSelector } from 'modules/auth/login/user.slice';
type Props = {
  onCloseModal: () => void;
};

const AddCollaboratorForm: FC<Props> = (props) => {
  const { onCloseModal } = props;
  const mutation = useCreateCollaborator();
  const toast = useToast();
  const { getAllCollaborators } = useFetchCollaborators();
  const { user } = useSelector(userSelector);
  const {
    createCollaborator,
    isCreating,
    wasCollaboratorCreated,
    hasCollaboratorFailed,
    errorMessage
  } = mutation;
  const {
    register,
    control,
    handleSubmit,
    formState: { errors }
  } = useForm<NewCollaborator>(ADD_COLLABORATOR_FORM_METHOD) as any;
  const [isChecked, setIsChecked] = useState(false);

  useEffect(() => {
    if (wasCollaboratorCreated) {
      getAllCollaborators();
    }
  }, [wasCollaboratorCreated, getAllCollaborators]);

  useEffect(() => {
    if (wasCollaboratorCreated) {
      onCloseModal();
      toast({
        position: 'top-right',
        title: t`Collaborator has been added successfully`,
        isClosable: true,
        duration: 9000,
        status: 'success'
      });
    }
  }, [wasCollaboratorCreated, onCloseModal, toast]);

  const onSubmit = (data: NewCollaborator) => {
    createCollaborator(data);
    // set collaborators in global state
  };

  useEffect(() => {
    if (hasCollaboratorFailed) {
        onCloseModal();
        toast({
            position: 'top-right',
            title: t`Unable to add collaborator`,
            description: t`${errorMessage}`,
            isClosable: true,
            duration: 9000,
            status: 'error'
        });
    }
  }, [hasCollaboratorFailed, onCloseModal, errorMessage, toast]);

  return (
    <chakra.form onSubmit={handleSubmit(onSubmit)}>
      <VStack mb={5}>
        <Text fontSize="sm" fontWeight={'bold'}>
          <Trans>Please Provide the Name and Email Address</Trans>
        </Text>
        <Text fontSize="sm">
          <Trans>
            The contact will receive an email to create a TRISA Global Directory Service (GDS)
            account or join this organization if the contact already has a TRISA GDS account. The
            invitation to join is valid for 7 calendar days. The contact will be added as a member
            for the VASP. The contact will have the ability to contribute to certificate requests,
            check on the status of certificate requests, and complete other actions related to the
            organization’s TRISA membership.
          </Trans>
        </Text>
      </VStack>

      {!isProdEnv ? <DevTool control={control} /> : null}
      <Stack>
        <Text fontWeight={'bold'} size={'md'}>
          <Trans>VASP</Trans>
        </Text>
        <Text data-testid="vasp-name">{user?.vasp?.name || 'No VASP name found'}</Text>
      </Stack>

      <Stack py={5}>
        <InputFormControl
          controlId="name"
          isInvalid={!!errors.name}
          data-test="name"
          formHelperText={errors.name?.message}
          {...register('name')}
          label={
            <>
              <chakra.span fontWeight={700}>
                <Trans>Contact Name</Trans>
              </chakra.span>{' '}
              (<Trans>required</Trans>)
            </>
          }
        />
      </Stack>
      <InputFormControl
        controlId="email"
        data-test="email"
        isInvalid={!!errors.email}
        formHelperText={errors?.email?.message}
        {...register('email')}
        label={
          <>
            <chakra.span fontWeight={700}>
              <Trans>Email Address</Trans>
            </chakra.span>{' '}
            (<Trans id="required">required</Trans>)
          </>
        }
      />

      <CheckboxFormControl
        controlId="agreed"
        mt={2}
        mb={4}
        {...register('agreed')}
        onChange={(e) => setIsChecked(e.target.checked)}
        size="md"
        borderColor={'black'}
        colorScheme="gray">
        <Trans>
          TRISA is a network of trusted members. I acknowledge that the contact is authorized to
          access the organization’s TRISA account information.
        </Trans>
      </CheckboxFormControl>

      <Stack display="flex" flexDir="column" gap={3} py={5}>
        <Button
          bg="orange"
          _hover={{ bg: 'orange' }}
          data-test="submit"
          minW="150px"
          type="submit"
          isDisabled={!isChecked}
          isLoading={isCreating}>
          Invite
        </Button>
        <Button
          variant="outline"
          mb={4}
          color="ghost"
          data-test="cancel"
          disabled={isCreating}
          fontWeight={400}
          onClick={onCloseModal}
          minW="150px">
          Cancel
        </Button>
      </Stack>
    </chakra.form>
  );
};

export default AddCollaboratorForm;
