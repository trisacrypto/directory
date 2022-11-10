
import React, { FC } from 'react';
import { useForm } from 'react-hook-form';
import { ADD_COLLABORATOR_FORM_METHOD } from './AddCollaboratorFormValidation';
import type { NewCollaborator } from './AddCollaboratorType';
import { useCreateCollaborator } from 'hooks/useCreateCollaborator';
import {
  Button,
  Stack,
  Text,
  VStack,
  chakra
} from '@chakra-ui/react';
import { DevTool } from '@hookform/devtools';
import { isProdEnv } from 'application/config';
import { Trans } from '@lingui/macro';
import InputFormControl from 'components/ui/InputFormControl';
import CustomToast from 'components/ui/CustomToast';
 import CheckboxFormControl from 'components/ui/CheckboxFormControl';

// import { setupI18n } from '@lingui/core';
// const _i18n = setupI18n();

type Props = {
    onCloseModal: () => void;
};

const AddCollaboratorForm: FC<Props> = (props) => {
    const { onCloseModal } = props;
    const mutation = useCreateCollaborator();

    const { createCollaborator, isLoading, isSuccess } = mutation;
    const {
        register,
        control,
        handleSubmit,
        formState: { errors },
    } = useForm(ADD_COLLABORATOR_FORM_METHOD) as any;

    const onSubmit = (data: NewCollaborator) => {
        createCollaborator(data);
    };

    if (isSuccess) {
        onCloseModal();
      return (
        <Stack data-test="collaborator-added">
          <CustomToast title="Collaborator added" status="success" />
        </Stack>
      );
    }

    return (

            <chakra.form onSubmit={handleSubmit(onSubmit)}>
                <VStack>
                    <Text>
                        Please provide the name of the VASP and email address for the new contact.
                    </Text>
                    <Text>
                        The contact will receive an email to create a TRISA Global Directory Service
                        Account. The invitation to join is valid for 7 calendar days. The contact will
                        be added as a member for the VASP. The contact will have the ability to
                        contribute to certificate requests, check on the status of certificate requests,
                        and complete other actions related to the organization’s TRISA membership.
                    </Text>
                </VStack>

                <CheckboxFormControl controlId="agreed" mt={2} mb={4} {...register('agreed')} colorScheme="green">
                    TRISA is a network of trusted members. I acknowledge that the contact is
                    authorized to access the organization’s TRISA account information.
                </CheckboxFormControl>


                {!isProdEnv ? <DevTool control={control} /> : null}
                <InputFormControl
                    controlId="name"
                    isInvalid={!!errors.name}
                    formHelperText={errors.name?.message}
                    {...register('name')}
                    label={
                        <>
                            <chakra.span fontWeight={700}>
                                <Trans>VASP Name</Trans>
                            </chakra.span>{' '}
                            (<Trans>required</Trans>)
                        </>
                    }
                />
                <InputFormControl
                    controlId="email"
                    isInvalid={!!errors.email}
                    formHelperText={errors?.email?.message}
                    {...register('email')}
                    label={
                        <>
                            <chakra.span fontWeight={700}>
                                <Trans>Email Address</Trans>
                            </chakra.span>{' '}
                            (<Trans>required</Trans>)
                        </>
                    }
                />

                <Stack display="flex" flexDir="column" gap={3}>
                    <Button
                        bg="orange"
                        _hover={{ bg: 'orange' }}
                        minW="150px"
                        type="submit"
                        isDisabled={isLoading}
                        isLoading={isLoading}>
                        Invite
                    </Button>
                    <Button
                        variant="ghost"
                        color="link"
                        disabled={isLoading}
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
