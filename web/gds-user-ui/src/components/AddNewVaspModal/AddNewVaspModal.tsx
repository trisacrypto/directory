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
import { canCreateOrganization } from 'utils/permission';
import AddNewVaspForm from '../AddNewVaspForm/AddNewVaspForm';

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

  // you dont' have permission to create a new organization

  return (
    <>
      <Button data-testid="add-new-vasp" onClick={onOpen} disabled={!canCreateOrganization()}>
        + Add New VASP
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
