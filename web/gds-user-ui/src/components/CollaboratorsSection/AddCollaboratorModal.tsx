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
  useDisclosure,
  VStack,
  chakra
} from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import { BsTrash } from 'react-icons/bs';
import InputFormControl from 'components/ui/InputFormControl';
function AddCollaboratorModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Link color="blue" onClick={onOpen}>
      <BsTrash fontSize="26px" />
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent w="100%" maxW="600px" px={10}>
          <ModalHeader textTransform="capitalize" textAlign="center" fontWeight={700} pb={1}>
            Add New Contact
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <VStack>
              <Text>
                Please provide the name of the VASP and email address for the new contact.
              </Text>
              <Text>
                The contact will receive an email to create a TRISA Global Directory Service
                Account. The invitation to join is valid for 7 calendar days. The contact will be
                added as a member for the VASP. The contact will have the ability to contribute to
                certificate requests, check on the status of certificate requests, and complete
                other actions related to the organization’s TRISA membership.
              </Text>
            </VStack>
            <Checkbox mt={2} mb={4}>
              TRISA is a network of trusted members. I acknowledge that the contact is authorized to
              access the organization’s TRISA account information.
            </Checkbox>

            <InputFormControl
              controlId="vasp_name"
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
            <Button bg="orange" _hover={{ bg: 'orange' }} minW="150px">
              Invite
            </Button>
            <Button variant="ghost" color="link" fontWeight={400} onClick={onClose} minW="150px">
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Link>
  );
}

export default AddCollaboratorModal;
