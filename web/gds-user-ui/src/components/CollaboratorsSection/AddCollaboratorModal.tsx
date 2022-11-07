/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useState, useEffect } from 'react';
import {
  Button,
  Checkbox,
  Link,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
    FormHelperText,
  FormErrorMessage,
  VStack,
  chakra
} from '@chakra-ui/react';
import { Trans, t } from '@lingui/macro';
import InputFormControl from 'components/ui/InputFormControl';
import { addCollaborator } from 'modules/dashboard/collaborator/service';
import { DevTool } from '@hookform/devtools';
import { isProdEnv } from 'application/config';
import { FormProvider, useForm, Controller } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { setupI18n } from '@lingui/core';
import * as yup from 'yup';
import CheckboxFormControl from 'components/ui/CheckboxFormControl';
const _i18n = setupI18n();
interface Props {
  isOpen: boolean;
  onClose: () => void;
  onOpen: () => void;
  onCloseModal: () => void;
}

function AddCollaboratorModal(props: Props) {
  const { isOpen, onOpen, onClose, onCloseModal } = props;
  const [isLoading, setIsLoading] = useState(false);
  const [_, setIsError] = useState(false);
  const resolver = yupResolver(
    yup
      .object()
      .shape({
        email: yup
          .string()
          .email(_i18n._(t`Email is not valid.`))
          .required(_i18n._(t`Email is required.`)),
        vasp_name: yup.string().required(),
        has_agreed: yup
          .boolean()
          .oneOf([true], _i18n._(t`You must agree to the terms and conditions`))
          .default(false)
      })
      .required()
  );

  const methods = useForm({
    defaultValues: {
      email: '',
      vasp_name: '',
      has_agreed: false
    },
    resolver,
    mode: 'onChange'
  });

  const {
    formState: { isDirty, errors },
    reset,
    register
  } = methods;

  console.log('[AddCollaboratorModal] errors', errors);
  function submitHandler() {
    console.log(methods.getValues());
    // setIsLoading(true);
    // try {
    //   const res = await addCollaborator({ ...methods.getValues() });
    //   setIsLoading(false);
    //   if (res.status === 200) {
    //     console.log(res);
    //   }
    //   console.log(res);
    // } catch (err) {
    //   setIsError(true);
    //   setIsLoading(false);
    // }
  }

  return (
    <Link color="blue" onClick={onOpen}>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent w="100%" maxW="600px" px={10}>
          <FormProvider {...methods}>
            <chakra.form onSubmit={methods.handleSubmit(submitHandler)}>
              <ModalHeader textTransform="capitalize" textAlign="center" fontWeight={700} pb={1}>
                Add New Contact
              </ModalHeader>
              <ModalCloseButton onClick={onCloseModal} />
              <ModalBody>
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

                <Checkbox mt={2} mb={4} {...register('has_agreed')} colorScheme="green">
                  TRISA is a network of trusted members. I acknowledge that the contact is
                  authorized to access the organization’s TRISA account information.
                </Checkbox>
                {!errors?.has_agreed ? (
                  <FormHelperText>{errors?.has_agreed}</FormHelperText>
                ) : (
                  <FormErrorMessage role="alert" data-testid="error-message">
                    {errors?.has_agreed}
                  </FormErrorMessage>
                )}

                {!isProdEnv ? <DevTool control={methods.control} /> : null}
                <InputFormControl
                  controlId="vasp_name"
                  isInvalid={!!errors.vasp_name}
                  formHelperText={errors.vasp_name?.message}
                  {...register('vasp_name')}
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
              </ModalBody>

              <ModalFooter display="flex" flexDir="column" gap={3}>
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
              </ModalFooter>
            </chakra.form>
          </FormProvider>
        </ModalContent>
      </Modal>
    </Link>
  );
}

export default AddCollaboratorModal;
