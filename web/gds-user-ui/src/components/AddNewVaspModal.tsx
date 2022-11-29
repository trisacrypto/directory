import {
  Button,
  Modal,
  ModalBody,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  useDisclosure,
  useToast
} from '@chakra-ui/react';
import { t, Trans } from '@lingui/macro';
import { usePostOrganizations } from 'modules/dashboard/organization/usePostOrganization';
import { FormProvider, useForm } from 'react-hook-form';
import * as Yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';
import { queryClient } from 'utils/react-query';
import { FETCH_ORGANIZATION } from 'constants/query-keys';
import AddNewVaspForm from './AddNewVaspForm';

const validationSchema = Yup.object().shape({
  name: Yup.string().required(t`The VASP Name is required.`),
  domain: Yup.string()
    .url(t`The Domain Name is invalid.`)
    .required(t`The Domain Name is required.`)
});

function AddNewVaspModal() {
  const { isOpen, onOpen, onClose: closeModal } = useDisclosure();
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
    reset,
    formState: { isSubmitting }
  } = methods;
  const { mutate, isLoading } = usePostOrganizations();
  const toast = useToast();
  const isCreatingVasp = isSubmitting || isLoading;

  const onSubmit = (values: any) => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { accept, ...payload } = values;

    mutate(payload, {
      onSuccess() {
        queryClient.invalidateQueries([FETCH_ORGANIZATION]);
        reset();
        closeModal();
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
      <Button onClick={onOpen}>
        + <Trans>Add New VASP</Trans>
      </Button>

      <Modal blockScrollOnMount isOpen={isOpen} onClose={closeModal}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader textAlign="center" textTransform="capitalize">
            <Trans>Add new managed VASP</Trans>
          </ModalHeader>

          <ModalBody>
            <FormProvider {...methods}>
              <AddNewVaspForm
                onSubmit={onSubmit}
                isCreatingVasp={isCreatingVasp}
                closeModal={closeModal}
              />
            </FormProvider>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}

export default AddNewVaspModal;
