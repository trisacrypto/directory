import { FC, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { ADD_COLLABORATOR_FORM_METHOD } from './AddCollaboratorFormValidation';
import type { NewCollaborator } from './AddCollaboratorType';
import { useCreateCollaborator } from 'hooks/useCreateCollaborator';
import { Button, Stack, Text, VStack, chakra, useToast } from '@chakra-ui/react';
import { DevTool } from '@hookform/devtools';
import { isProdEnv } from 'application/config';
import { Trans } from '@lingui/react';
import InputFormControl from 'components/ui/InputFormControl';
import CustomToast from 'components/ui/CustomToast';
import CheckboxFormControl from 'components/ui/CheckboxFormControl';
import { t } from '@lingui/macro';
import { useFetchCollaborators } from 'components/Collaborators/useFetchCollaborator';
type Props = {
  onCloseModal: () => void;
};

const AddCollaboratorForm: FC<Props> = (props) => {
  const { onCloseModal } = props;
  const mutation = useCreateCollaborator();
  const toast = useToast();
  const { getAllCollaborators } = useFetchCollaborators();

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
    getValues: formValues,
    formState: { errors }
  } = useForm<NewCollaborator>(ADD_COLLABORATOR_FORM_METHOD) as any;

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

  if (hasCollaboratorFailed) {
    return (
      <Stack data-test="collaborator-failed">
        <CustomToast
          title="Fail to add a new collaborator"
          description={errorMessage}
          status="error"
        />
      </Stack>
    );
  }

  return (
    <chakra.form onSubmit={handleSubmit(onSubmit)}>
      <VStack>
        <Text>
          <Trans id="Please provide the name of the VASP and email address for the new contact.">
            Please provide the name of the VASP and email address for the new contact.
          </Trans>
        </Text>
        <Text>
          <Trans id="The contact will receive an email to create a TRISA Global Directory Service Account. The invitation to join is valid for 7 calendar days. The contact will be added as a member for the VASP. The contact will have the ability to contribute to certificate requests, check on the status of certificate requests, and complete other actions related to the organization’s TRISA membership.">
            The contact will receive an email to create a TRISA Global Directory Service Account.
            The invitation to join is valid for 7 calendar days. The contact will be added as a
            member for the VASP. The contact will have the ability to contribute to certificate
            requests, check on the status of certificate requests, and complete other actions
            related to the organization’s TRISA membership.
          </Trans>
        </Text>
      </VStack>

      <CheckboxFormControl
        controlId="agreed"
        mt={2}
        mb={4}
        data-test="agreed"
        {...register('agreed')}
        colorScheme="gray">
        <Trans id="TRISA is a network of trusted members. I acknowledge that the contact is authorized to access the organization’s TRISA account information.">
          TRISA is a network of trusted members. I acknowledge that the contact is authorized to
          access the organization’s TRISA account information.
        </Trans>
      </CheckboxFormControl>

      {!isProdEnv ? <DevTool control={control} /> : null}
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
                <Trans id="VASP Name">VASP Name</Trans>
              </chakra.span>{' '}
              (<Trans id="required">required</Trans>)
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
              <Trans id="Email Address">Email Address</Trans>
            </chakra.span>{' '}
            (<Trans id="required">required</Trans>)
          </>
        }
      />

      <Stack display="flex" flexDir="column" gap={3} py={5}>
        <Button
          bg="orange"
          _hover={{ bg: 'orange' }}
          data-test="submit"
          minW="150px"
          type="submit"
          isDisabled={formValues().agreed === false}
          isLoading={isCreating}>
          Invite
        </Button>
        <Button
          variant="outline"
          mb={4}
          color="link"
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
