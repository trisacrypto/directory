import {
  Button,
  chakra,
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
import { Trans } from '@lingui/macro';
import UnlinkIcon from 'assets/copy-link.svg';
import InputFormControl from 'components/ui/InputFormControl';

function RemoveLinkedAccountModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <IconButton
        variant="unstyled"
        aria-label="copy the link"
        icon={<CkLazyLoadImage src={UnlinkIcon} mx="auto" w="25px" />}
        marginTop="25px!important"
        onClick={onOpen}
      />

      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent border="1px solid black" px={5}>
          <ModalHeader textAlign="center" fontWeight={700} fontSize="md">
            <Trans>Remove Linked Account</Trans>
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Trans>
              The following account will be removed from your profile. The linked account will not
              be affected.
            </Trans>

            <VStack>
              <InputFormControl
                controlId="linked_account"
                label={
                  <chakra.span fontWeight={700} textTransform="capitalize" mt={5} display="block">
                    Linked Account
                  </chakra.span>
                }
                value="jones.ferdinand@luminous.co.uk"
              />
            </VStack>
          </ModalBody>

          <ModalFooter display="flex" flexDir="column" rowGap={2}>
            <Button bg="orange" _hover={{ bg: 'orange' }} minW="150px">
              <Trans>Confirm</Trans>
            </Button>
            <Button variant="ghost" onClick={onClose}>
              <Trans>Cancel</Trans>
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

export default RemoveLinkedAccountModal;
