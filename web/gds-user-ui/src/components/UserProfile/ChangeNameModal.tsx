import {
  Button,
  FormLabel,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  useDisclosure,
  VStack
} from '@chakra-ui/react';
import CkLazyLoadImage from 'components/LazyImage';
import EditIcon from 'assets/edit-input.svg';
import { Trans } from '@lingui/macro';
import InputFormControl from 'components/ui/InputFormControl';

function ChangeNameModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <IconButton
        aria-label="Edit"
        icon={<CkLazyLoadImage src={EditIcon} mx="auto" w="25px" />}
        variant="unstyled"
        marginTop="32px!important"
        onClick={onOpen}
      />
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent border="1px solid black" px={5}>
          <ModalHeader textAlign="center" fontWeight={700} fontSize="md">
            Change Name
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <VStack align="start">
              <InputFormControl
                label={
                  <FormLabel fontWeight={700}>
                    <Trans>First (Given) Name</Trans>
                  </FormLabel>
                }
                controlId="first_given_name"
              />
              <InputFormControl
                label={
                  <FormLabel fontWeight={700}>
                    <Trans>Last (Family) Name</Trans>
                  </FormLabel>
                }
                controlId="first_given_name"
              />
            </VStack>
          </ModalBody>

          <ModalFooter display="flex" flexDir="column" rowGap={2}>
            <Button bg="orange" _hover={{ bg: 'orange' }} minW="150px">
              Save
            </Button>
            <Button variant="ghost" onClick={onClose}>
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

export default ChangeNameModal;
