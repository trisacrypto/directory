import {
  Button,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  Text,
  useDisclosure,
  chakra,
  useToast
} from '@chakra-ui/react';
import { t, Trans } from '@lingui/macro';
import { usePostOrganizations } from 'modules/dashboard/organization/usePostOrganization';
import { FormProvider, useForm } from 'react-hook-form';
import CheckboxFormControl from './ui/CheckboxFormControl';
import InputFormControl from './ui/InputFormControl';
import * as Yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';
import { queryClient } from 'utils/react-query';
import { FETCH_ORGANIZATION } from 'constants/query-keys';

const validationSchema = Yup.object().shape({
  name: Yup.string().required(t`The VASP Name is required.`),
  domain: Yup.string()
    .url(t`The Domain Name is invalid.`)
    .required(t`The Domain Name is required.`)
});

function AddNewVaspModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const methods = useForm({
    defaultValues: {
      name: '',
      domain: '',
      accept: false
    },
    mode: 'onSubmit',
    resolver: yupResolver(validationSchema)
  });
  const {
    register,
    watch,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting }
  } = methods;
  const { mutate, isLoading } = usePostOrganizations();
  const toast = useToast();
  const isCreatingVasp = isSubmitting || isLoading;

  const accept = watch('accept');
  const onSubmit = (values: any) => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { accept: _accept, ...payload } = values;

    mutate(payload, {
      onSuccess() {
        queryClient.invalidateQueries([FETCH_ORGANIZATION]);
        reset();
        onClose();
      },
      onError: (error) => {
        console.log('[mutate] error', error.response?.data.error);
        toast({
          title: error.response?.data?.error || error.message,
          status: 'error',
          position: 'top-right'
        });
      }
    });
  };

  return (
    <>
      <Button onClick={onOpen}>+ Add New VASP</Button>

      <Modal blockScrollOnMount isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader textAlign="center">Modal Title</ModalHeader>

          <ModalBody>
            <FormProvider {...methods}>
              <form onSubmit={handleSubmit(onSubmit)}>
                <Text>
                  <Trans>
                    Please input the name of the new managed Virtual Asset Service Provider (VASP).
                    When the entity is created, you will have the ability to add collaborators,
                    start and complete the certificate registration process, and manage the VASP
                    account. Please acknowledge below and provide the name of the entity.
                  </Trans>
                </Text>
                <CheckboxFormControl
                  controlId="accept"
                  {...register('accept', { required: true })}
                  colorScheme="gray">
                  <Trans>
                    TRISA is a network of trusted members. I acknowledge that the new VASP has a
                    legitimate business purpose to join TRISA.
                  </Trans>
                </CheckboxFormControl>
                <Stack py={5}>
                  <InputFormControl
                    controlId="name"
                    isInvalid={!!errors.name}
                    data-testid="name"
                    formHelperText={errors.name?.message}
                    isDisabled={!accept}
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
                    controlId="domain"
                    isInvalid={!!errors.domain}
                    data-testid="domain"
                    formHelperText={errors.domain?.message}
                    isDisabled={!accept}
                    placeholder="https://"
                    {...register('domain')}
                    label={
                      <>
                        <chakra.span fontWeight={700}>
                          <Trans>VASP Domain</Trans>
                        </chakra.span>{' '}
                        (<Trans>required</Trans>)
                      </>
                    }
                  />
                </Stack>
                <ModalFooter display="flex" flexDir="column" justifyContent="center" gap={2}>
                  <Button
                    bg="orange"
                    _hover={{ bg: 'orange' }}
                    type="submit"
                    minW={150}
                    isDisabled={!accept || isCreatingVasp}>
                    <Trans id="Next">Create</Trans>
                  </Button>
                  <Button variant="ghost" onClick={onClose} disabled={isCreatingVasp}>
                    <Trans id="Cancel">Cancel</Trans>
                  </Button>
                </ModalFooter>
              </form>
            </FormProvider>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}

export default AddNewVaspModal;
